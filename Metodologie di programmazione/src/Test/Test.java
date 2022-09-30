package Test;

import java.io.File;
import java.nio.file.FileSystem;
import java.nio.file.Path;
import java.util.Arrays;

public class Test
{
    public static void main(String[] args)
    {
        String[] dettagli = new String[]{"ciao", "come", "va", "tutto", "bene"};
        dettagli = Arrays.copyOfRange(dettagli, 1, dettagli.length);
        System.out.println(Arrays.toString(dettagli));
    }

}
