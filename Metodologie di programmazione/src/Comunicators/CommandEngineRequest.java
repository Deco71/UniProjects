package Comunicators;

public class CommandEngineRequest
{
    private String method;
    private String[] classes;
    private String[] miscelaneous;

    CommandEngineRequest(String method, String[] classes, String[] miscelaneous)
    {
        this.method = method;
        this.classes = classes;
        this.miscelaneous = miscelaneous;
    }

    CommandEngineRequest(String method)
    {
        this(method, null, null);
    }
    CommandEngineRequest(String method, String[] classes) { this(method, classes, null); }

    public String getMethod()
    {
        return method;
    }

    public String[] getClasses()
    {
        return classes;
    }

    public String[] getMiscelaneous()
    {
        return miscelaneous;
    }
}
