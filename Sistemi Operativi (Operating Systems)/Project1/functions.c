#include <stdio.h>
#include <stdlib.h>
#include <fcntl.h>
#include <string.h>
#include <stdbool.h>
#include <unistd.h> //necessario per avere dup2() su sistemi linux, non necessario su sistemi windows

// ---------------------------------------------- //
// ----------- FUNCTIONS DECLARATIONS ----------- //
// ---------------------------------------------- //

// Questa funzione restituisce una struct vuota e inizializzata utile per salvare le informazioni degli argomenti
struct Options get_options();

/* Questa funzione aggiorna la struct in base all'argomento passatogli in input
 * struct Options: Il riferimento alla struct di tipo Options da aggiornare
 * char **: Il riferimento ad argv
 * int: Il numero massimo di argomenti presenti in argv
 * int: L'argomento che dobbiamo analizzare
 * RETURNS: Il prossimo argomento i che dobbiamo analizzare
 */
int set_options(struct Options *, char **, int, int);

/* Questa funzione imposta l'output su file
 * struct Options: Il riferimento alla struct di tipo Options con le informazioni necessarie
 */
void set_output(struct  Options *);

/* Questa funzione legge 50 righe per volta del file fornito in input e le inserisce all'interno dell'array fornito
 * FILE: Il puntatore al file pointer da analizzare
 * char **: Il puntatore all'array da aggiornare con i nuovi dati
 * RETURNS: Il numero di righe lette
 */
int read_file_rows(FILE *, char **);

/* Questa funzione ci analizza due stringhe e ci comunica se sono uguali o diverse
 * char * e char *: Le due stringhe
 * RETURNS: True se le due stringhe sono uguali, False se sono diverse
 */
bool check_rows(char *, char *);

/* Questa funzione libera la memoria in fase di chiusura del programma
 * int * e int *: Rispettivamente le righe lette del primo e del secondo file
 * char ** e char**: Rispettivamente i due array con i puntatori alle righe dei file
 */

// ---------------------------------------------- //
// ------------- STRUCT DECLARATIONS ------------ //
// ---------------------------------------------- //

struct Options {
    //Nomi dei file da aprire
    char *first_file; //Nome del primo file da aprire
    char *second_file; //Nome del secondo file da aprire

    //Variabili per l'output
    bool option_o; //Impostato a True se l'argomento -o è settato, fa il redirect dell'output sul file
    char default_output[10]; //Nome del file di default se -o è impostato ma non è stato fornito nessun nome
    char *output_file; //Nome del file di output fornito dall'utente

    //Opzioni combinate
    bool option_g; //Impostato a True se l'argomento -g è settato, restituisce una stringa quando i file sono diversi
    bool option_s; //Impostato a True se l'argomento -s è settato, restituisce una stringa quando i file sono uguali

    //Opzioni singole (Ma possono essere usate insieme alle opzioni combinate)
    bool option_d; //Impostato a True se l'argomento -d è settato, restituisce il numero delle righe diverse
    bool option_u; //Impostato a True se l'argomento -u è settato, restituisce il numero delle righe uguali
    //Print delle righe (da usare solo insieme a una opzione singola)
    bool option_v; //Impostato a True se l'argomento -v è settato, imposta la modalità verbose del programma
};


// ---------------------------------------------- //
// ----------- FUNCTIONS DEFINITIONS ------------ //
// ---------------------------------------------- //

struct Options get_options()
{
    struct Options options;
    options.option_o = false;
    options.first_file = NULL;
    options.second_file = NULL;
    options.output_file = NULL;
    strcpy(options.default_output, "./out.txt");
    options.option_g = false;
    options.option_s = false;
    options.option_d = false;
    options.option_u = false;
    options.option_v = false;
    return options;
}

int set_options(struct Options* options, char **argv, int argc, int i)
{
    //Il valore della variabile i preso in input che dovremo restituire alla fine
    int return_i = i;
    //La lunghezza dell'argomento letto
    int len = strlen(argv[i]);
    //Iteriamo la lunghezza dell'argomento partendo da uno (escludendo quindi -)
    for (int j = 1; j < len; j++)
    {
        //Prendiamo l'argomento e lo confrontiamo con tutte le possibili opzioni...
        char argument = argv[i][j];
        if(argument == 'o')
        {
            //se l'argomento successivo in argv (se esiste) non è -, allora è il nome del file di output
            if ((i+1) < argc && argv[i+1][0] != '-')
            {
                options->option_o = true;
                options->output_file = argv[i+1];
                //inoltre dobbiamo dire al main che il prossimo argomento non va controllato
                return_i = i+1;
            }
                //altrimenti mettiamo quello di default
            else
            {
                options->option_o = true;
                options->output_file = options->default_output;
            }
            //A questo punto chiamiamo set_output
            set_output(options);
        }
        //Stessa cosa per più o meno tutti gli argomenti
        if(argument == 'g')
        {
            options->option_g = true;
        }
        if(argument == 's')
        {
            options->option_s = true;
        }
        if(argument == 'd')
        {
            //se u è già impostato, restituisci errore
            if(options->option_u)
                return -1;
            options->option_d = true;
        }
        if(argument == 'u')
        {
            //se d è già impostato, restituisci errore
            if(options->option_d)
                return -1;
            options->option_u = true;
        }
        if(argument == 'v')
        {
            options->option_v = true;
        }
    }
    return return_i;
}

void set_output(struct Options *options)
{
    //Facciamo il flush dello standard output (buona pratica prima di chiuderlo e/o duplicarlo)
    fflush(stdout);
    //Apriamo il file di output (creandolo se necessario)
    int fd1 = open(options->output_file, O_WRONLY | O_CREAT | O_TRUNC, 0644);
    //Se abbiamo un problema nell'aprire/creare il file, abortiamo l'operazione
    if (fd1 < 0)
    {
        perror("Errore nell'apertura del file\n");
        exit(-1);
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

int read_file_rows(FILE *fp, char **pointer_buffer)
{
    char *buffer;
    char ch;
    int i;
    //altrimenti iteriamo per 50 volte (o fino a quando non arriviamo alla fine del file)
    for (i = 0; i < 50; i++)
    {
        //Ok contiene false se la riga non è ancora stata letto per intero, true altrimenti
        bool ok = false;
        //Grandezza iniziale del buffer
        int Buffer_Dim = 1024;
        //Controlliamo se la malloc ha avuto esito positivo
        if ((buffer = (char *) malloc(Buffer_Dim)) == NULL)
        {
            perror("Errore nella allocazione della memoria");
            exit(-1);
        }
        int length;
        //Fino a quando la riga non abbiamo letto tutta la riga...
        while(!ok)
        {
            //...prendi la riga e la sua lunghezza
            fgets(buffer, Buffer_Dim, fp);
            length = strlen(buffer);
            //Controlliamo l'ultimo elemento della riga, se non è \n e non siamo arrivati alla fine del file...
            ch = buffer[length-1];
            if(ch != '\n' && !feof(fp))
            {
                //...aumenta il buffer...
                Buffer_Dim *= 2;
                //...libera la memoria precedentemente allocata...
                free(buffer);
                //...rialloca la memoria (non utilizziamo realloc() perchè non vogliamo i dati della vecchia allocazione)
                if ((buffer = malloc(Buffer_Dim)) == NULL)
                {
                    perror("Errore nella allocazione della memoria");
                    exit(-1);
                }
                //Spostiamo poi il cursore all'inizio della riga e ritentiamo la lettura
                fseek(fp, -length, SEEK_CUR);
            }
            else
                ok = true;
        }
        //Quando abbiamo letto tutta la riga allochiamo la memoria nel pointer buffer
        if ((pointer_buffer[i] = malloc(Buffer_Dim)) == NULL)
        {
            perror("Errore nella allocazione della memoria");
            exit(-1);
        }
        //copiamo il contenuto della memoria puntata da buffer all'interno della memoria puntata da pointer buffer
        memcpy(pointer_buffer[i], buffer, Buffer_Dim);
        //liberiamo la memoria allocata in buffer
        free(buffer);
        //se il file è finito prima delle 50 righe, usciamo dal ciclo
        if (feof(fp))
        {
            break;
        }
    }
    //restituiamo il numero di righe lette
    return i;
}


bool check_rows(char * row1, char * row2)
{
    if (strcmp(row1, row2) == 0)
    {
        return true;
    } else
        return false;
}
