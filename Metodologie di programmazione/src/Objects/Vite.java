package Objects;

import Comunicators.Condition;
import Comunicators.Creator;
import Comunicators.MethodResponse;
import Comunicators.Phrase;

import java.util.List;
import java.util.stream.Collectors;

public class Vite extends Oggetto implements Prendibile
{
    private List<String> usabileSu;

    Vite(Creator c, String s)
    {
        super(c, s);
        usabileSu = c.getObjectList().stream().map(String::valueOf).collect(Collectors.toList());
    }

    @Override
    public MethodResponse guarda(Condition c)
    {
        return new MethodResponse(Phrase.VitLook);
    }

    @Override
    public MethodResponse apri(Condition c)
    {
        return new MethodResponse(Phrase.VitOpen);
    }

    @Override
    public MethodResponse rompi(Condition c)
    {
        return baseStringNo;
    }

    @Override
    public MethodResponse usa(Condition c)
    {
        return new MethodResponse(Phrase.VitUse);
    }

    @Override
    public MethodResponse prendi(Condition c)
    {
        c.getPersonaggio().setInventario(this);
        return baseStringPick;
    }
}
