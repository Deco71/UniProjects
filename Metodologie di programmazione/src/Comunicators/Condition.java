package Comunicators;

import Mondo.Room;
import Objects.Oggetto;
import Personaggi.Personaggio;

import java.lang.management.BufferPoolMXBean;
import java.util.ArrayList;

public class Condition
{
    private Personaggio personaggio;
    private ArrayList<Oggetto> oggetti;
    public Room stanza;

    Condition(Personaggio p, Room stanza, ArrayList<Oggetto> o)
    {
        this.personaggio = p;
        this.stanza = stanza;
        this.oggetti = o;
    }

    public Personaggio getPersonaggio()
    {
        return personaggio;
    }

    public ArrayList<Oggetto> getOggetti()
    {
        return oggetti;
    }

    public Oggetto getOggetti(int index)
    {
        if(oggetti.size() > index+1)
            return oggetti.get(index);
        return null;
    }

    public Room getStanza()
    {
        return stanza;
    }
}
