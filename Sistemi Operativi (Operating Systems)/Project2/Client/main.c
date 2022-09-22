#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <sys/socket.h>
#include <math.h>
#include <arpa/inet.h>
#include <unistd.h>

//Struct con all'interno le informazioni sulla connessione in corso
struct connection_info
{
    int socket_fd;
    struct sockaddr_in server_addr;
};

/*
 * Funzione che restituisce una struct di tipo connection_info che tenta ad oltranza la connessione
 * con il server fino ad esito positivo
*/
struct connection_info get_connection();


int main()
{
    //Cerchiamo di stabilire una connessione
    struct connection_info connessione = get_connection();

    //Appena otteniamo una connessione stabile con il server, printiamo alcune informazioni sull'applicativo
    printf("Connesso!\n");
    printf("FORMATO DELLE ISTRUZIONI\n");
    printf("Numero; un carattere tra +, -, *, /; Numero\n");
    printf("L'applicazione prende anche i numeri negativi e con la virgola\n");
    printf("Per chiudere l'applicazione, scrivere il comando exit\n");
    printf(">: ");
    fflush(stdout);

    //-------VARIABILI--------//

    //Array con le informazioni che invieremo al server per l'elaborazione
    double invio[3];

    //Un valore booleano che sarà impostato a true quando il client rileva una perdita di connessione con il server
    int connection_lost = 0;

    //Un intero che conterrà il numero di elementi correttamente parsati dalla scanf
    int n;

    //Il veccchio timestamp dello scambio di infomazioni precedenti (servirà per effettuare un controllo sullo stato della socket)
    double old_timestamp = -1;

    //Intero che ci indica se abbiamo avuto un errore in fase di lettura dell'input e che tipo di errore
    //1 = L'operatore inserito non è valido
    //2 = E' stata tentata una divisione per 0
    int error = 0;

    //Carattere che conterrà il nostro "presunto" operatore
    char operatore;

    //Buffer che pulirà il contenuto di stdin dopo la lettura dell'input
    char cleaner[1024];

    //Array che conterrà gli elementi ricevuti dal server elaborati
    double ricevuti[3];

    //Ciclo infinito (dobbiamo soddisfare tante richieste quante volute dall'utente)
    while(1)
    {
        /*
        * Se nell'iterazione precedente non abbiamo perso la connessione, chiedi l'operazione che l'utente vuole inviare al server
        * altrimenti, invia semplicemente l'operazione, senza chiedere all'utente di reimmetterla
        */
        if (connection_lost == 0)
        {
            n = 0;
            error = 0;
            n = scanf("%lf %c %lf", &invio[0], &operatore, &invio[2]);

            //Puliamo lo stdin da eventuali residui
            int i = 0;
            while (1) {
                scanf("%c", &cleaner[i]);
                if (cleaner[i] == '\n')
                    break;
                i++;
            }
            cleaner[i] = '\0';
        }
        //Se non abbiamo preso correttamente tutti e 3 i campi, o ci sono dei residui, avvisa l'utente...
        if (n != 3 || strcmp("", cleaner) != 0)
        {
            //Se lo scanf iniziale ha fallito e la stringa rimasta nel buffer era proprio exit, usciamo dal programma.
            if (n != 3 && strcmp("exit", cleaner) == 0)
            {
                break;
            }
            printf("Errore nell'inserimeno dei dati, attenersi alla struttura specificata all'inizio del programma\n");
            printf(">: ");
            fflush(stdout);
        }
        //...altrimenti analizza la stringa in input e vedi se il carattere inserito è un operatore valido
        else
        {
            connection_lost = 0;
            switch (operatore)
            {
                case '+':
                    invio[1] = 0;
                    break;
                case '-':
                    invio[1] = 1;
                    break;
                case '*':
                    invio[1] = 2;
                    break;
                case '/':
                    //Se l'utente prova a fare una divisione per 0, annulla tutto
                    if (invio[2] == 0)
                        error = 2;
                    invio[1] = 3;
                    break;
                //Se l'operatore non è nessuno dei tre, annulla tutto 
                default:
                    printf("Inserisci un operatore valido\n");
                    printf(">: ");
                    fflush(stdout);
                    error = 1;
                    break;
            }
            if (error == 2)
            {
                printf("Non puoi dividere per zero!\n");
                printf(">: ");
            }
            //Se non abbiamo errori, procediamo con l'invio delle informazioni al server
            else if (error == 0)
            {
                if ((write(connessione.socket_fd, invio, sizeof(invio))) < 0)
                {
                    perror("Errore mentre scrivevo sul socket");
                    exit(EXIT_FAILURE);
                }
                printf("In attesa di risposta\n");
                fflush(stdout);
                read(connessione.socket_fd, ricevuti, sizeof(ricevuti));
                /*
                 * Workaround per cercare di aggiungere resilienza al client
                 * Se dalla socket riceviamo un vecchio timestamp o dati che non hanno senso,
                 * allora probabilmente il server è andato offline.
                 * Chiudiamo quindi il socket attuale e cerchiamo d'instaurare una
                 * nuova connessione pulita.
                 * Una volta ottenuta la nuova connessione, ritentare l'invio dell'operazione.
                 */
                if (old_timestamp == ricevuti[0] || fabs(ricevuti[0]) < 0.1)
                {
                    close(connessione.socket_fd);
                    printf("Connessione con il server interrotta\n");
                    printf("Stiamo ritentando la connessione, attendere... \n");
                    connection_lost = 1;
                    fflush(stdout);
                    connessione = get_connection();
                    printf("Connessione stabilita nuovamente!\n");
                }
                //Se il risultato ricevuto sembra pulito, allora stampiamo il risultato a schermo
                else
                {
                    old_timestamp = ricevuti[0];
                    double service_time = ricevuti[1] - ricevuti[0];
                    printf("Risultato: %lf (%lf secondi di elaborazione)\n", ricevuti[2], service_time);
                    printf(">: ");
                    fflush(stdout);
                }

            }
        }
    }
    //Se usciamo dal programma, printiamo un saluto e chiudiamo il socket
    printf("Arrivederci\n");
    close(connessione.socket_fd);
    return 0;
}

struct connection_info get_connection()
{
    //Definiamo una struct di tipo connection_info
    struct connection_info connessione;

    //Inizializziamo la socket
    connessione.socket_fd = socket(AF_INET, SOCK_STREAM, 0);
    if (connessione.socket_fd < 0)
    {
        perror("Errore nella creazione del socket");
        exit(EXIT_FAILURE);
    }

    //Assegnamo alla socket le informazioni necessarie a connettersi al server
    connessione.server_addr.sin_family = AF_INET;
    connessione.server_addr.sin_port = htons(7131);
    //Per semplicita didattica ho hardcodato l'indirizzo di loopback, ma nulla ci vieta di inserire all'interno un IP reale preso da stdin
    inet_pton(AF_INET, "127.0.0.1", &connessione.server_addr.sin_addr);

    //Avvisiamo l'utente che un tentativo di connessione al server è in corso
    printf("Tentando una connesione... attendere\n");
    fflush(stdout);

    //Finche non riusciamo a stabilire una connessione, tentiamo; quando riusciamo restituiamo la struct di tipo connection_info
    while(connect(connessione.socket_fd, (struct sockaddr *)&connessione.server_addr, sizeof(connessione.server_addr)) < 0)
    {
        sleep(3);
    }
    return connessione;
}