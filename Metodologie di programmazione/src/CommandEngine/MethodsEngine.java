package CommandEngine;

import Comunicators.CommandEngineRequest;
import Personaggi.Personaggio;

import java.io.BufferedReader;
import java.io.IOException;
import java.lang.reflect.Method;
import java.util.ArrayList;
import java.util.HashMap;

public class MethodsEngine extends Engine
{
    private HashMap<SupportedLanguages, HashMap<String, Method>> methodsMap = new HashMap<>();
    private String section = "Methods";

    protected MethodsEngine()
    {
        super.engineStarter(section);
    }

    @Override
    public void engineSetup(BufferedReader br, SupportedLanguages lingua) throws IOException, NoSuchMethodException
    {
        ArrayList<String> baseArray = super.baseFileLoader(section);
        HashMap<String, Method> map = new HashMap<>();
        boolean sectionOn = false;
        int instructionsInserted = 0;
        while(br.ready())
        {
            String line = br.readLine();
            if(line.contains("["))
            {
                System.out.println("entro");
                sectionOn = findSection(line, section);
                System.out.println(sectionOn);
            }
            else if(sectionOn && !line.equals(""))
            {
                String[] dettagli = line.split(" ");
                String nomeMetodo = baseArray.get(instructionsInserted);
                Method m = Personaggio.class.getMethod(nomeMetodo, CommandEngineRequest.class);
                for(String s : dettagli)
                {
                    map.put(s, m);
                }
                instructionsInserted++;
            }
        }
        if(instructionsInserted != baseArray.size())
        {
            super.error(section, instructionsInserted, baseArray.size());
        }
        System.out.println(map);
        methodsMap.put(lingua, map);
        System.out.println(methodsMap + " Methods");
    }

}
