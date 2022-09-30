import CommandEngine.LanguageEngine;

import java.util.HashMap;

/**
 * La classe Action racchiude il nome dei comandi che dobbiamo chiamare utilizzando la reflection.
 * E' un tassello essenziale del {@link LanguageEngine}
 */
public class Action
{
    /**
     * Contiene il nome del metodo da chiamare sottoforma di stringa
     */
    private String nomemetodo;
    /**
     * Contiene la mappa da stringa a azione per salvarsi tutte le azioni create.
     * Necessaria per implementare il singleton pattern
     */
    private static HashMap<String, Action> mappaAzioni = new HashMap<>();

    /**
     * Costruttore privato per rispettare il singleton pattern
     * @param nomemetodo Il nome del metodo che l'azione dovrà richiamare
     */
    private Action(String nomemetodo)
    {
        this.nomemetodo = nomemetodo;
    }

    /**
     * Restituisce il nome del metodo
     * @return Il nome del metodo
     */
    public String getNomemetodo()
    {
        return nomemetodo;
    }

    /**
     * Crea/restituisce un oggetto di tipo Action
     * @param nomeclasse Il nome della classe
     * @return
     */
    public static Action getAction(String nomeclasse)
    {
        mappaAzioni.computeIfAbsent(nomeclasse, Action::new);
        return mappaAzioni.get(nomeclasse);
    }

    public static void main(String[] args)
    {
        Action a = getAction("Saluto");
        Action b = getAction("Saluto");
        System.out.println(a == b);
        System.out.println(a.getNomemetodo());
        System.out.println(b.getNomemetodo());
    }
}
