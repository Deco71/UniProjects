#include <stdio.h>
#include <stdlib.h>
#include <stdbool.h>
#include <errno.h>
#include "functions.c"

// ---------------------------------------------- //
// ------------------ MAIN ---------------------- //
// ---------------------------------------------- //
int main(int argc, char **argv)
{
    // ------------------ Inizializzazione e fetch degli argomenti ---------------------- //

    //Prendiamo la struct di tipo options inizializzata
    struct Options options;
    options = get_options();


    //Iteriamo poi per ogni argomento presente nell'input
    for (int i = 1; i < argc; i++)
    {
        if (argv[i][0] == '-')
        {
            //se il primo elemento dell'argomento è - allora dobbiamo vedere se è una opzione
            i = set_options(&options, argv, argc, i);
        }
        else
        {
            //altrimenti sarà o il primo file...
            if (!options.first_file)
            {
                options.first_file = argv[i];
            }
            //o il secondo
            else if (!options.second_file)
            {
                options.second_file = argv[i];
            }
            //nel caso ci sia un terzo file, interrompere l'esecuzione
            else
            {
                i = -1;
            }
        }
        //Nel caso ci sia stato un errore con un qualche argomento non valido, restituiamo un errore e usciamo
        if (i==-1)
        {
            //Impostiamo errno a 22 (Invalid Argument)
            errno = 22;
            perror("Argomenti non validi");
            exit(EXIT_FAILURE);
        }
    }
    //Se non è specificato nessun argomento di rilevanza, esci dal programma
    if (!options.option_g && !options.option_s && !options.option_d && !options.option_u)
    {
        errno = 22;
        perror("Argomenti non validi");
        exit(EXIT_FAILURE);
    }
    //Usciamo inoltre anche se l'utente non ha fornito i nomi di due file
    if(!options.second_file || !options.first_file)
    {
        errno = 22;
        perror("Definisci due file, tu ne ha definiti uno o zero");
        exit(EXIT_FAILURE);
    }

    // ------------------ Apertura del file e operazioni di controllo ---------------------- //

    //Apriamo il primo file e il secondo
    FILE *fp1 = fopen(options.first_file, "r");
    FILE *fp2 = fopen(options.second_file, "r");
    //Se uno dei due file pointer è NULL, c'è stato un problema con l'apertura dei file
    if (fp1 == NULL || fp2 == NULL)
    {
        perror("Errore nell'apertura dei file da leggere");
        exit(EXIT_FAILURE);
    }

    /* Allochiamo i nostri puntatori ai pointer buffer, che saranno coloro che conterranno i puntatori alle locazioni
     * di memoria dove saranno situate le stringhe dei file
    */
    char * pointer_buffer_1[50];
    char * pointer_buffer_2[50];

    //Il numero di righe lette per ogni file
    int rows1;
    int rows2;

    //Il numero generale di righe lette
    int total_rows = 0;

    //Contiene il risultato della operazione di controllo delle stringhe
    bool result;

    //Contiene true se per il momento i file sono uguali, false altrimenti
    bool different = false;

    //Contiene il numero dell'iterazione (ci servirà per liberare la memoria)
    int i;

    //Iteriamo fino a quando uno dei due file non finisce
    while (( !(feof(fp2)) && !(feof(fp1)) ))
    {
        //Ci prendiamo le nostre 50 righe (o meno nel caso il file sia finito)
        rows1 = read_file_rows(fp1, pointer_buffer_1);
        rows2 = read_file_rows(fp2, pointer_buffer_2);
        //Li confrontiamo riga per riga
        for (i = 0; i < rows1 && i < rows2; i++)
        {
            //Sommiamo al numero totale di righe la riga appena letta (partiamo da uno per la prima riga)
            total_rows ++;
            //Controlliamo l'uguaglianza delle righe
            result = check_rows(pointer_buffer_1[i], pointer_buffer_2[i]);
            //Se le righe sono diverse...
            if (!result)
            {
                //...impostiamo different a true
                different = true;
                //Inoltre se g è anche l'unica opzione settata (non dobbiamo fare un controllo riga per riga)
                if(!options.option_d && !options.option_u && options.option_g)
                {
                    //Printiamo la seguente stringa e usciamo dal programma
                    printf("I file %s e %s sono differenti\n", options.first_file, options.second_file);
                    return 0;
                }
            }
            //Se -d o -u è settato e ci sono le condizioni appropriate...
            if((options.option_d && !result) || (options.option_u && result))
            {
                //...printiamo il numero della stringa uguale/diversa...
                printf("%d", total_rows);
                //...e se è settato anche -v...
                if(options.option_v)
                {
                    //...printiamo anche le stringhe (modalità verbosa)
                    if((feof(fp1) || feof(fp2)) && ( i == rows1-1 || i == rows2-1))
                    {
                        //se siamo alla fine dei file, mandiamo a capo per rendere il tutto più leggibile
                        printf(", %s: %s\n", options.first_file, pointer_buffer_1[i]);
                        printf("%d, %s: %s\n", total_rows, options.second_file, pointer_buffer_2[i]);
                    }
                    else
                    {
                        printf(", %s: %s", options.first_file, pointer_buffer_1[i]);
                        printf("%d, %s: %s", total_rows, options.second_file, pointer_buffer_2[i]);
                    }
                }
                printf("\n");
                //E poi liberiamo dalla memoria le due righe già lette
                free(pointer_buffer_1[i]);
                free(pointer_buffer_2[i]);
            }
        }
    }
    //Visto che stiamo per uscire dal programma, prima di tutto liberiamo la memoria rimasta allocata
    for (int j = i + 1; j < rows1; j++)
    {
        free(pointer_buffer_1[j]);
    }
    for (int j = i + 1; j < rows2; j++)
    {
        free(pointer_buffer_2[j]);
    }
    //Se i file hanno lunghezze diverse (un file è finito prima di un'altro)...
    if ( ( ( feof(fp1) || feof(fp2) ) && rows1 != rows2 ))
    {
        //...e -g è settata, riportare che sono diversi...
        if (options.option_g)
        {
            printf("I file %s e %s sono differenti\n", options.first_file, options.second_file);
        }
        //...altrimenti ritorna senza dire nulla
        return 0;
    }
    //...invece se alla fine del ciclo i file non sono risultati diversi e -s è settato, printa questa stringa...
    if (!different && options.option_s)
    {
        printf("I file %s e %s sono uguali\n", options.first_file, options.second_file);
    }
    else if (different && options.option_g) //...altrimenti se -g è settato printa questa.
        printf("I file %s e %s sono differenti\n", options.first_file, options.second_file);
    //Concludiamo
    return 0;
}
