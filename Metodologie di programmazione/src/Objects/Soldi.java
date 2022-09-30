package Objects;

import Comunicators.Condition;
import Comunicators.Creator;
import Comunicators.MethodResponse;
import Comunicators.Phrase;

public class Soldi extends Oggetto implements Prendibile
{

    Soldi(Creator c)
    {
        super(c, "Soldi");
    }

    @Override
    public MethodResponse guarda(Condition c)
    {
        return new MethodResponse(Phrase.SolLook);
    }

    @Override
    public MethodResponse apri(Condition c)
    {
        return baseStringNo;
    }

    @Override
    public MethodResponse rompi(Condition c)
    {
        return baseStringNo;
    }

    @Override
    public MethodResponse usa(Condition c)
    {
        return new MethodResponse(Phrase.SolUse);
    }

    @Override
    public MethodResponse prendi(Condition c)
    {
        c.getPersonaggio().setInventario(this);
        return new MethodResponse(Phrase.SolGet);
    }
}
