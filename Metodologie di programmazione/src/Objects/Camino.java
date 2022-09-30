package Objects;

import Comunicators.Condition;
import Comunicators.Creator;
import Comunicators.MethodResponse;

public class Camino extends Oggetto implements Accendible
{

    Camino(Creator c)
    {
        super(c, "Camino");
    }

    boolean status = false;

    @Override
    public boolean isOn()
    {
        return false;
    }

    @Override
    public void on()
    {

    }

    @Override
    public void off()
    {

    }

    public boolean On()
    {
        return true;
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
