package Test;

public class FuckAndShit
{
    /*
    private static Engine.LingueSupportate linguaAttuale = Engine.LingueSupportate.ITA;
    private static String directory = "LanguageFiles";
    private static String baseFile = "BaseFile.txt";
    private static Engine.LanguageEngine instance;
    private static HashMap<Engine.LingueSupportate, HashMap<String, Method>> methodsMap = new HashMap<>();
    private static HashMap<Engine.LingueSupportate, HashMap<String, Class<?>>> classMap = new HashMap<>();
    private static HashMap<Engine.LingueSupportate, HashMap<String, String>> namesMap = new HashMap<>();

    private LocalizationEngine()
    {
    }

    public Requester.LocalizationResponse localizator(Requester.LocalizationRequest request)
    {
        if(request.toString().toLowerCase().contains("change"))
        {
            String[] s = request.toString().split(" ");
            return new Requester.LocalizationResponse(changeLanguage(s[s.length-1]), true);
        }
        else
        {
            //String s = methodsMap.get(linguaAttuale).getOrDefault(request.getMethod(), null);
            if(null == null)
                return new Requester.LocalizationResponse(linguaAttuale.getError(), true);
            return null;
        }
    }

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


    public static Engine.LanguageEngine localizatorCreator()
    {
        if(instance == null)
        {
            try
            {
                instance = localizationStarter();
            }
            catch (NotEnoughtInstructionsException e)
            {
                System.out.println("Huston we have a 'probliem' - LocalizationEngine localizationSettuper Error");
                e.printStackTrace();
                System.exit(-1);
            }
        }
        linguaAttuale = Engine.LingueSupportate.ITA;
        return instance;
    }

    private static Engine.LanguageEngine localizationStarter() throws NotEnoughtInstructionsException
    {
        Engine.LanguageEngine engine = new Engine.LanguageEngine();
        ArrayList<String> baseArray = baseFileLoader();
        Path file = Path.of(directory);
        File[] files = file.toFile().listFiles();
        for(File f : files)
        {
            int instructionsInserted = 0;
            if(f.getName().equals(baseFile))
                continue;
            HashMap<String, String> languageMap = new HashMap<>();
            try(BufferedReader br = new BufferedReader(new FileReader(f)))
            {
                while(br.ready())
                {
                    String linea = br.readLine();
                    String[] dettagli = linea.split(" ");
                    String nomeClasse = baseArray.get(instructionsInserted);
                    for(String s : dettagli)
                    {
                        languageMap.put(s, nomeClasse);
                    }
                    instructionsInserted++;
                }
                if(instructionsInserted != baseArray.size())
                {
                    System.out.println(instructionsInserted);
                    System.out.println(baseArray.size());
                    throw new NotEnoughtInstructionsException();
                }
            }
            catch (IOException e)
            {
                System.out.println("Huston we have a 'probliem' - LocalizationEngine Starter Error");
                e.printStackTrace();
                System.exit(-1);
            }
            methodsMap.put(findEnum(f.getName()), languageMap);
        }
        return engine;
    }

    private static ArrayList<String> baseFileLoader()
    {
        ArrayList<String> baseFileArray = new ArrayList<>();
        try(BufferedReader br = new BufferedReader(new FileReader(Paths.get(directory, baseFile).toString())))
        {
            while(br.ready())
                baseFileArray.add(br.readLine());
        }
        catch (IOException e)
        {
            System.out.println("Huston we have a 'probliem' - LocalizationEngine baseFileLoader Error");
            e.printStackTrace();
            System.exit(-1);
        }
        return baseFileArray;
    }

    private static Engine.LingueSupportate findEnum(String lingua)
    {
        String[] l = lingua.split("//.");
        return Enum.valueOf(Engine.LingueSupportate.class, l[0].substring(0, 3).toUpperCase());
    }
}

class NotEnoughtInstructionsException extends Throwable
{

}
     */

}
