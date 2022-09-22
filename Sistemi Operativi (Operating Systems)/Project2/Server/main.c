#include <arpa/inet.h>
#include <signal.h>
#include <time.h>
#include <stdlib.h>
#include <stdio.h>
#include <pthread.h>
#include <semaphore.h>
#include <string.h>
#include <sys/socket.h>
#include <unistd.h>
#include <fcntl.h>

//Struttura della quale passeremo l'indirizzo ai thread che creiamo
struct thread_args
{
    //L'id del thread
    pthread_t id;
    //Il socket con il client
    int client_sock;
    //L'IP del client
    char* client_ip;
};

//Il semaforo che ci permetterà di gestire correttamente il flusso di scrittura sul log
pthread_mutex_t print_semaphore;

/*
 * Questa funzione prende come input l'indirizzo di memoria della struct di tipo thread_args con gli argomenti del thread
 * Inoltre, verrà passato il puntatore ad essa in fase di creazione del thread, visto che contiene le istruzioni che esso
 * dovrà effettuare
 */
void *operazioni(void*);

// Questa funzione imposta l'output su file
void set_output();

//Questa funzione restituisce un timestamp in formato double
double get_timestamp();


int main()
{
    //----------VARIABILI---------//
    int server_socket;
    struct sockaddr_in server_info, clientAddr;
    sigset_t set;

    //Inseriamo all'interno della struct sockaddr_in le informazioni per inizializzare la socket TCP/IP//
    server_info.sin_addr.s_addr = htonl(INADDR_ANY); //Prendiamo tutti gli indirizzi della macchina
    server_info.sin_family = AF_INET; //Impostiamo la famiglia di indirizzi IPV4
    server_info.sin_port = htons(7131); //Impostiamo una porta

    //Inizializziamo la struttura dati socket creando una unnamed socket e impostando il protocollo TCP/IP (SOCK_STREAM)
    server_socket = socket(AF_INET, SOCK_STREAM, 0);
    if (server_socket == -1)
    {
        perror("Errore nella apertura del socket");
        exit(EXIT_FAILURE);
    }

    //Associamo gli indirizzi specificati in server_info alla socket facendo il bind
    if(bind(server_socket, (struct sockaddr *) &server_info, sizeof(server_info)) < 0 )
    {
        perror("Errore nel binding del socket");
        exit(EXIT_FAILURE);
    }

    //Inizializziamo l'ascolto della socket, impostando il numero massimo di connessioni in attesa a 10
    if(listen(server_socket, 10) < 0)
    {
        perror("Errore nell'inizializzazione dell'ascolto");
        exit(EXIT_FAILURE);
    }

    //Impostiamo l'output su file del log delle operazioni
    set_output();

    //Impostiamo il nostro semaforo necessario per garantire la mutua esclusione sulla scrittura sul log
    pthread_mutex_init(&print_semaphore, NULL);

    //Aggiungiamo al set di signal il segnale SIGCHLD e poi gli aggiungiamo il SIG_BLOCK per bloccarlo
    sigemptyset(&set);
    sigaddset(&set, SIGCHLD);
    sigprocmask(SIG_BLOCK, &set, NULL);
    //Così facendo, non dovremo fare la wait() sui nostri thread morenti e non diverranno zombie

    //Iniziamo il ciclo infinito per accettare le connessioni fino a terminazione forzata del server
    while(1)
    {
        socklen_t cli_len = sizeof(clientAddr);
        int new_sock;
        //Accettiamo la nuova connessione in ingresso inserendo in new_sock il file descriptor della socket
        new_sock = accept(server_socket, (struct sockaddr *) &clientAddr, &cli_len);
        if (new_sock < 0)
        {
            perror("Errore nella accettazione della connessione");
        }
        else
        {
            //Nel caso la connessione abbia esito positivo, ci salviamo in formato stringa l'IP del client
            char ipClient[INET_ADDRSTRLEN];
            inet_ntop(AF_INET, &clientAddr.sin_addr, ipClient, INET_ADDRSTRLEN);

            //Definiamo poi thread_arguments...
            struct thread_args *thread_arguments;
            //...che prepareremo prima allocando lo spazio necessario con una malloc...
            thread_arguments = malloc(sizeof(struct thread_args));
            //...e poi passandogli le informazioni necessarie al corretto funzionamento del thread
            thread_arguments->client_sock = new_sock;
            thread_arguments->client_ip = ipClient;

            /*
             * Infine, lanciamo il nuovo thread con la start_routine impostata alla funzione operazioni()
             * e come argomenti l'indirizzo di thread_arguments che libereremo alla fine del thread.
             */
            pthread_create(&thread_arguments->id, NULL, operazioni, (void*) thread_arguments);
        }
    }
}


void *operazioni(void *args)
{
    //Prendiamo e facciamo il cast al puntatore degli argomenti
    struct thread_args *arguments;
    arguments = (struct thread_args *) args;
    //Cerchiamo di entrare nella regione critica per stampare le informazioni con l'avvenuta connessione
    pthread_mutex_lock(&print_semaphore);
    printf("Connessione accettata da %s\n", arguments->client_ip);
    fflush(stdout);
    //Liberiamo la regione critica
    pthread_mutex_unlock(&print_semaphore);

    //E ci prepariamo a ricevere la richiesta del client
    double ricezione[3];
    char operazione;
    double invio[3];
    //Finché il client è connesso e pronto a inviare dati...
    while(read(arguments->client_sock, ricezione, sizeof(ricezione)))
    {
        //Prendiamo il timestamp della ricezione
        invio[0] = get_timestamp();
        //Controlliamo il tipo di operazioni richiesta dal client
        switch ((int) ricezione[1])
        {
            case 0:
                invio[2] = ricezione[0]+ricezione[2];
                operazione = '+';
                break;
            case 1:
                invio[2] = ricezione[0]-ricezione[2];
                operazione = '-';
                break;
            case 2:
                invio[2] = ricezione[0]*ricezione[2];
                operazione = '*';
                break;
            case 3:
                invio[2] = ricezione[0]/ricezione[2];
                operazione = '/';
                break;
            default:
                invio[2] = 0;
                operazione = '?';
                break;
        }
        //Prendiamo il timestamp di invio del messaggio e lo inviamo subito dopo
        invio[1] = get_timestamp();
        if ((write(arguments->client_sock, invio, sizeof(invio))) < 0)
        {
            perror("Errore mentre scrivevo sul socket");
        }
        else //Se l'invio delle informazioni avviene in modo corretto...
        {
            //...richiedere l'accesso alla regione critica e inserire le informazioni sull'operazione svolta nel log
            pthread_mutex_lock(&print_semaphore);
            printf("IP: %s; Operazione: %c; Risultato: %lf; Arrivo Richiesta: %lf; Restuituzione Richiesta: %lf\n",
                   arguments->client_ip, operazione, invio[2], invio[0], invio[1]);
            fflush(stdout);
            pthread_mutex_unlock(&print_semaphore);
        }

    }
    //Se la connessione cade, prendiamo entriamo nella regione critica...
    pthread_mutex_lock(&print_semaphore);
    //...notifichiamo della perdita di connessione con il client...
    printf("La connessione con %s è caduta\n", arguments->client_ip);
    fflush(stdout);
    //...e rilasciamo nuovamente il semaforo
    pthread_mutex_unlock(&print_semaphore);
    //Infine, rimuoviamo dalla memoria gli argomenti relativi al thread
    close(arguments->client_sock);
    free(arguments);
    return 0;
}

void set_output()
{
    //Facciamo il flush dello standard output (buona pratica prima di chiuderlo e/o duplicarlo)
    fflush(stdout);
    //Apriamo il file di output (creandolo se necessario)
    int fd1 = open("log.txt", O_WRONLY | O_CREAT | O_TRUNC, 0644);
    //Se abbiamo un problema nell'aprire/creare il file, abortiamo l'operazione
    if (fd1 < 0)
    {
        perror("Errore nell'apertura del file\n");
        exit(EXIT_FAILURE);
    }
    //altrimenti "inseriamo" il nostro file descriptor all'interno di quello dello standard output
    if (dup2(fd1, STDOUT_FILENO) < 0)
    {
        printf("Errore nella duplicazione del file descriptor\n");
        exit(-1);
    }
    //e chiudiamo il file
    close(fd1);
}


/*
 * LA NOSTRA FUNZIONE E' THREAD SAFE?
 * Clock_gettime() è thread safe
 * sprintf() è thread safe, a meno che il buffer non è condiviso tra i thread, cosa che nel nostro caso non accade
 * strtod() è thread safe (sul man è segnata come MT-safe locale)
 *
 * La nostra funzione è quindi thread safe da usare per i nostri scopi visto che non modifichiamo la locale
 * durante l'esecuzione
 */
double get_timestamp()
{
    struct timespec spec;
    char time_str[32];

    //Prendiamo il tempo attuale dell'orologio del PC (tempo calcolato a partire dall'EPOC)
    clock_gettime(CLOCK_REALTIME, &spec);
    //Lo trasformiamo in stringa
    sprintf(time_str, "%ld.%.3ld", spec.tv_sec, spec.tv_nsec);
    //E poi lo trasformiamo in double
    double time = strtod(time_str, NULL);
    return time;
}