package Objects;

import Comunicators.Condition;
import Comunicators.Creator;
import Comunicators.MethodResponse;
import Comunicators.Phrase;

public class Secchio extends Oggetto implements Accendible {
    Secchio(Creator c, String nomeOggetto)
    {
        super(c, nomeOggetto);
    }

    boolean status = false;

    @Override
    public MethodResponse guarda(Condition c)
    {
        return new MethodResponse(Phrase.VitOpen);
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

    @Override
    public boolean isOn()
    {
        return status;
    }

    @Override
    public void on()
    {
        status = true;
    }

    @Override
    public void off()
    {
        status = false;
    }
}
