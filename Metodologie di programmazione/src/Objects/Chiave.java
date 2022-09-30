package Objects;

import Comunicators.Condition;
import Comunicators.Creator;
import Comunicators.MethodResponse;

public class Chiave extends Oggetto
{

    Chiave(Creator c, String nomeOggetto)
    {
        super(c, nomeOggetto);
    }

    @Override
    public MethodResponse guarda(Condition c)
    {
        return null;
    }

    @Override
    public MethodResponse apri(Condition c)
    {
        return null;
    }

    @Override
    public MethodResponse rompi(Condition c)
    {
        return null;
    }

    @Override
    public MethodResponse usa(Condition c)
    {
        return null;
    }
}
