package Comunicators;

public class CommandEngineResponse
{
    private String method;
    private Class<?>[] classes;
    private String[] miscelaneous;
    private boolean isMessage = false;

    public CommandEngineResponse()
    {

    }

    CommandEngineResponse(String method, Class<?>[] classes, String[] miscelaneous)
    {
        this.method = method;
        this.classes = classes;
        this.miscelaneous = miscelaneous;
    }

    CommandEngineResponse(String message, boolean isMessage)
    {
        this(message, null, null);
        this.isMessage = isMessage;
    }

    public String getMethod()
    {
        return method;
    }

    public Class<?>[] getClasses()
    {
        return classes;
    }

    public String[] getMiscelaneous()
    {
        return miscelaneous;
    }

    public boolean isMessage()
    {
        return isMessage;
    }
}
