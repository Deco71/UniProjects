package CommandEngine;

import java.io.BufferedReader;
import java.io.IOException;
import java.util.ArrayList;
import java.util.HashMap;

public class ClassesEngine extends Engine
{
    private HashMap<SupportedLanguages, HashMap<String, Class<?>>> classesMap = new HashMap<>();
    private String section = "Classes";


    protected ClassesEngine()
    {
        super.engineStarter(section);
    }


    @Override
    public void engineSetup(BufferedReader br, SupportedLanguages lingua) throws IOException, NoSuchMethodException, ClassNotFoundException
    {
        ArrayList<String> baseArray = super.baseFileLoader(section);
        HashMap<String, Class<?>> map = new HashMap<>();
        boolean sectionOn = false;
        int instructionsInserted = 0;
        while(br.ready())
        {
            String line = br.readLine();
            if(line.contains("["))
            {
                sectionOn = findSection(line, section);
            }
            else if(sectionOn && !line.equals(""))
            {
                String nomeClasse = baseArray.get(instructionsInserted);
                Class<?> c = Class.forName("Elements."+nomeClasse);
                map.put(line, c);
                instructionsInserted++;
            }
        }
        if(instructionsInserted != baseArray.size())
        {
            super.error(section, instructionsInserted, baseArray.size());
        }
        classesMap.put(lingua, map);
        System.out.println(classesMap + " Class");
    }
}
