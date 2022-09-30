package CommandEngine;

public class EngineFactory
{
    private ClassesEngine instanceOfClasses = null;
    private LanguageEngine instanceOfLanguage = null;
    private MethodsEngine instanceOfMethods = null;

    public ClassesEngine startClasses()
    {
        if(instanceOfClasses == null)
            instanceOfClasses = new ClassesEngine();
        return instanceOfClasses;
    }

    public LanguageEngine startLanguage()
    {
        if(instanceOfLanguage == null)
            instanceOfLanguage = new LanguageEngine();
        return instanceOfLanguage;
    }

    public MethodsEngine startMethods()
    {
        if(instanceOfMethods == null)
            instanceOfMethods = new MethodsEngine();
        return instanceOfMethods;
    }


}
