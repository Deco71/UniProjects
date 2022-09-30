package Personaggi;

import Comunicators.CommandEngineRequest;
import Comunicators.CommandEngineResponse;
import Objects.Oggetto;

import java.util.ArrayList;

public class Personaggio
{
    public String nome;
    public ArrayList<Oggetto> inventario;

    public CommandEngineResponse guarda(CommandEngineRequest r)
    {
        return new CommandEngineResponse();
    }

    public CommandEngineResponse apri(CommandEngineRequest r)
    {
        return new CommandEngineResponse();
    }

    public CommandEngineResponse vai(CommandEngineRequest r)
    {
        return new CommandEngineResponse();
    }

    public CommandEngineResponse prendi(CommandEngineRequest r)
    {
        return new CommandEngineResponse();
    }

    public CommandEngineResponse rompi(CommandEngineRequest r)
    {
        return new CommandEngineResponse();
    }

    public CommandEngineResponse usa(CommandEngineRequest r)
    {
        return new CommandEngineResponse();
    }

    public CommandEngineResponse inventario(CommandEngineRequest r)
    {
        return new CommandEngineResponse();
    }

    public CommandEngineResponse dai(CommandEngineRequest r)
    {
        return new CommandEngineResponse();
    }

    public boolean inInventario(Oggetto o)
    {
        return inventario.indexOf(o) > 0;
    }
    public void setInventario(Oggetto o)
    {
        inventario.add(o);
    }
}
