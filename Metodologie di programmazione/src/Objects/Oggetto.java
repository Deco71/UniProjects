package Objects;

import Comunicators.Condition;
import Comunicators.Creator;
import Comunicators.MethodResponse;
import Comunicators.Phrase;
import Mondo.Room;

public abstract class Oggetto
{
    protected static final MethodResponse baseStringNo = new MethodResponse(Phrase.BaseNo);
    protected static final MethodResponse baseStringPick = new MethodResponse(Phrase.BasePick);
    protected static final MethodResponse baseStringBreakNo = new MethodResponse(Phrase.BaseBreakNo);
    protected static final MethodResponse baseStringAlreadyBroke = new MethodResponse(Phrase.BaseAlreadyBroke);

    private String nomeOggetto;
    private boolean rotto = false;
    protected Creator creator;
    //protected Room stanza;

    Oggetto(Creator c, String nomeOggetto)
    {
        this.nomeOggetto = nomeOggetto;
        this.creator = c;
    }


    public abstract MethodResponse guarda(Condition c);
    public abstract MethodResponse apri(Condition c);
    public abstract MethodResponse rompi(Condition c);
    public abstract MethodResponse usa(Condition c);

    /*
    //TODO: Implementare il controllo dell'usable in fase di chiamata dell'oggetto altrimenti son cazzi
    public boolean isUsable(Condition c)
    {
        return c.getStanza().getRoomName().equals(creator.getRoom().getRoomName()) || c.getPersonaggio().inInventario(this);
    }*/

    public String getNomeOggetto()
    {
        return nomeOggetto;
    }

    public boolean isRotto()
    {
        return rotto;
    }

    public void setRotto()
    {
        rotto = true;
    }

    public void aggiusta()
    {
        rotto = false;
    }

}
