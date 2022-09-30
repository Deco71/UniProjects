package Gioco;

import CommandEngine.ClassesEngine;
import CommandEngine.EngineFactory;
import CommandEngine.LanguageEngine;
import CommandEngine.MethodsEngine;
import Mondo.Mondo;

import java.util.Scanner;

public class Gioco
{
    private EngineFactory e = new EngineFactory();
    private LanguageEngine langEngine = e.startLanguage();
    private ClassesEngine classEngine = e.startClasses();
    private MethodsEngine methEngine = e.startMethods();

    public void play(Mondo m)
    {
        System.out.println("Benvenuto nel gioco! Inizia scrivendo un comando!\nIf you want you can change language using the command change with the name of the language you want to switch");
        Scanner s = new Scanner(System.in);
        while(true)
        {
            String comando = s.nextLine();
        }
    }
}
