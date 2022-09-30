package Objects;

import Comunicators.Condition;
import Comunicators.Creator;
import Comunicators.MethodResponse;
import Comunicators.Phrase;

import java.util.ArrayList;
import java.util.List;

public class Scrivania extends Oggetto
{
    private ArrayList<Oggetto> o = new ArrayList<>();
    private boolean open = false;

    Scrivania(Creator c, String s)
    {
        super(c,s);
        o.addAll(c.getObjectList());
    }

    @Override
    public MethodResponse guarda(Condition c)
    {
        if(open)
            return new MethodResponse(Phrase.ScrLook, o.toString().replace("[", "").replace("]", ""));
        return new MethodResponse(Phrase.ScrLookNo);
    }

    @Override
    public MethodResponse apri(Condition c)
    {
        open = true;
        return new MethodResponse(Phrase.ScrOpen, o);
    }

    @Override
    public MethodResponse rompi(Condition c)
    {
        return baseStringBreakNo;
    }

    @Override
    public MethodResponse usa(Condition c)
    {
        return apri(c);
    }
}
