package Comunicators;

import Mondo.Room;
import Objects.Oggetto;

import java.util.ArrayList;
import java.util.List;

public class Creator
{
    private ArrayList<Oggetto> objectList = new ArrayList<>();
    private Room m;
    Creator(Room m, Oggetto... o)
    {
        objectList.addAll(List.of(o));
        this.m = m;
    }

    public ArrayList<Oggetto> getObjectList()
    {
        return objectList;
    }

    public Oggetto getObjectList(int i)
    {
        return objectList.get(i);
    }

    public Room getRoom()
    {
        return m;
    }
}
