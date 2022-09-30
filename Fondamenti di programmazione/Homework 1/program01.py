def ex1(g1, g2, g3, g4, dim_hand, mazzo):
    lettere = [g1, g2, g3, g4]
    punti = [0, 0, 0, 0]
    carte = [dim_hand, dim_hand, dim_hand, dim_hand]
    dizio = {'e': 1, 'a': 1, 'i': 1, 'o': 1, 'n': 1, 'r': 1, 't': 1, 'l': 1, 's': 1, 'u': 1, 'b': 3, 'c': 3, 'm': 3,'p': 3, 'f': 4, 'h': 4, 'v': 4, 'w': 4, 'y': 4, 'd': 2, 'g': 2, 'q': 10, 'z': 10, 'j': 8, 'x': 8, 'k': 5}
    range0 = range(0, 4)
    mazzo -= (dim_hand * 4)
    for indice in range(len(g1)):
        for turn in range0:
            parola = lettere[turn][indice]
            mazzo -= len(parola)
            punti[turn] += sum([dizio[lettera] for lettera in parola])
            if mazzo <= 0:
                carte[turn] += mazzo
                mazzo = 0
                if carte[turn] == 0:
                    primo = (turn+1) % 4
                    secondo = (turn+2) % 4
                    terzo = (turn+3) % 4
                    punti[primo] -= (carte[primo] * 3)
                    punti[secondo] -= (carte[secondo] * 3)
                    punti[terzo] -= (carte[terzo] * 3)
                    return punti
