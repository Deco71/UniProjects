package Mondo;

import CommandEngine.ClassesEngine;
import CommandEngine.EngineFactory;
import CommandEngine.LanguageEngine;
import CommandEngine.MethodsEngine;

public class Mondo
{
    private EngineFactory e = new EngineFactory();
    private LanguageEngine langEngine = e.startLanguage();
    private ClassesEngine classEngine = e.startClasses();
    private MethodsEngine methEngine = e.startMethods();
    private Mondo()
    {

    }
    public static Mondo fromFile(String nomeFile)
    {

        return new Mondo();
    }
}
