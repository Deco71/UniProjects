package CommandEngine;

import java.io.BufferedReader;
import java.io.IOException;
import java.util.ArrayList;
import java.util.HashMap;

public class LanguageEngine extends Engine
{
    private HashMap<SupportedLanguages, HashMap<String, String>> languageMap = new HashMap<>();
    private String section = "Language";


    protected LanguageEngine()
    {
        super.engineStarter(section);
    }


    @Override
    public void engineSetup(BufferedReader br, SupportedLanguages lingua) throws IOException, NoSuchMethodException
    {
        ArrayList<String> baseArray = super.baseFileLoader(section);
        HashMap<String, String> map = new HashMap<>();
        boolean sectionOn = false;
        int instructionsInserted = 0;
        while(br.ready())
        {
            String line = br.readLine();
            if(line.contains("["))
            {
                sectionOn = findSection(line, section);
            }
            else if(sectionOn && baseArray.size() > instructionsInserted)
            { ;
                String traduzione = baseArray.get(instructionsInserted);
                map.put(line, traduzione);
                instructionsInserted++;
            }
        }
        if(instructionsInserted != baseArray.size())
        {
            super.error(section, instructionsInserted, baseArray.size());
        }
        languageMap.put(lingua, map);
        System.out.println(languageMap + " Language");
    }
}
