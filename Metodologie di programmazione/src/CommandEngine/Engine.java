package CommandEngine;

import java.io.*;
import java.nio.file.Path;
import java.nio.file.Paths;
import java.util.ArrayList;

public abstract class Engine
{

    protected String directory = "MiniZak";
    protected String baseFile = "BaseFile.txt";
    protected SupportedLanguages linguaAttuale = SupportedLanguages.ITA;

    public abstract void engineSetup(BufferedReader br, SupportedLanguages l) throws IOException, NoSuchMethodException, ClassNotFoundException;

    public String changeLanguage(String linguanuova)
    {
        try
        {
            linguaAttuale = findEnum(linguanuova);
        }
        catch (IllegalArgumentException e)
        {
            return "We still doesn't support that language, we are deeply saddened...";
        }
        return linguaAttuale.getMess();
    }

    private static SupportedLanguages findEnum(String lingua)
    {
        String[] l = lingua.split("//.");
        return Enum.valueOf(SupportedLanguages.class, l[0].substring(0, 3).toUpperCase());
    }

    protected void engineStarter(String section)
    {
        Path file = Path.of(directory);
        File[] files = file.toFile().listFiles();
        for(File f : files)
        {
            boolean sectionOn = false;
            if(f.getName().equals(baseFile))
                continue;
            try(BufferedReader br = new BufferedReader(new FileReader(f)))
            {
                engineSetup(br, findEnum(f.getName()));
            }
            catch (IOException | NoSuchMethodException | ClassNotFoundException e)
            {
                System.out.println("Huston we have a 'probliem' - LocalizationEngine Starter Error");
                e.printStackTrace();
                System.exit(-1);
            }
        }
    }

    protected ArrayList<String> baseFileLoader(String section)
    {
        ArrayList<String> baseFileArray = new ArrayList<>();
        try(BufferedReader br = new BufferedReader(new FileReader(Paths.get(directory, baseFile).toString())))
        {
            boolean sectionOn = false;
            while(br.ready())
            {
                String line = br.readLine();
                if(line.contains("["))
                {
                    sectionOn = findSection(line, section);
                }
                else if(sectionOn && !line.equals(""))
                    baseFileArray.add(line);
            }
        }
        catch (IOException e)
        {
            System.out.println("Huston we have a 'probliem' - LocalizationEngine baseFileLoader Error");
            e.printStackTrace();
            System.exit(-1);
        }
        return baseFileArray;
    }

    protected boolean findSection(String line, String section)
    {
        line = line.replace("[", "").replace("]", "");
        return line.equals(section);
    }

    protected void error(String dove, int instructionsInserted, int baseArray)
    {
        System.out.println(instructionsInserted);
        System.out.println(baseArray);
        System.out.println("Errore in " + dove);
        System.exit(-1);
    }
}