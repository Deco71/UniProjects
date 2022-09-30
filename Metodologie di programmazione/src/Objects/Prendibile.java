package Objects;

import Comunicators.Condition;
import Comunicators.MethodResponse;

public interface Prendibile
{
    MethodResponse prendi(Condition c);
}
