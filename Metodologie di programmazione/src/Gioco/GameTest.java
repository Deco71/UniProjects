package Gioco;

import Mondo.Mondo;

public class GameTest
{
    public static void main(String[] args)
    {
        Gioco g = new Gioco();
        Mondo m = Mondo.fromFile("hi");
        g.play(m);
    }
}
