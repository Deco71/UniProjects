package CommandEngine;

public enum SupportedLanguages
{
    ITA {public String getMess(){ return "La tua lingua da adesso in poi sarà l'italiano, per cambiarla utilizza il comando change";}
         public String getError(){ return "Non ho capito cosa vuoi fare...";}},
    ENG {public String getMess(){ return "Your language is now English, you can change that using the change command";}
         public String getError(){ return "I didn't understand what you want to do...";}};

    public abstract String getMess();
    public abstract String getError();
}
