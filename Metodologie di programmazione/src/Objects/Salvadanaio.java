package Objects;

import Comunicators.Condition;
import Comunicators.Creator;
import Comunicators.MethodResponse;
import Comunicators.Phrase;

import java.util.ArrayList;
import java.util.List;
import java.util.Objects;

public class Salvadanaio extends Oggetto
{
    private Oggetto contiene;
    private Oggetto rottoDa;

    Salvadanaio(Creator c, String s)
    {
        super(c,s);
        contiene = creator.getObjectList(0);
        rottoDa = creator.getObjectList(1);
    }

    @Override
    public MethodResponse guarda(Condition c)
    {
        return new MethodResponse(Phrase.SalLook);
    }

    @Override
    public MethodResponse apri(Condition c)
    {
        if(isRotto())
        {
            return baseStringAlreadyBroke;
        }
        if(c.getOggetti().get(0).getNomeOggetto().equals(rottoDa.getNomeOggetto()))
        {
            setRotto();
            return new MethodResponse(Phrase.SalBreakYes, List.of(contiene));
        }
        return baseStringBreakNo;
    }

    @Override
    public MethodResponse rompi(Condition c)
    {
        return apri(c);
    }

    @Override
    public MethodResponse usa(Condition c)
    {
        return apri(c);
    }
}
