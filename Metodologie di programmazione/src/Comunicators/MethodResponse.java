package Comunicators;

import Objects.Oggetto;

import java.util.List;

public class MethodResponse
{
    private Phrase phrase;
    private String string;
    private List<Oggetto> newObject;

    public MethodResponse(Phrase p, String s, List<Oggetto> o)
    {
        this.phrase = p;
        this.string = s;
        this.newObject = o;
    }

    public MethodResponse(Phrase p)
    {
        this(p, null, null);
    }

    public MethodResponse(Phrase p, String s)
    {
        this(p, s, null);
    }

    public MethodResponse(Phrase p, List<Oggetto> o)
    {
        this(p, null, o);
    }

    public Phrase getPhrase()
    {
        return phrase;
    }

    public String getString()
    {
        return string;
    }

    public List<Oggetto> getNewObject()
    {
        return newObject;
    }
}
