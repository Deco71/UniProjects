def ex1(istructions_file, initial_city, clues):
    clues = clues.split()
    risultato = recursive(trasforma(creazionelista(istructions_file)), initial_city, 0, clues, "", set(), clues.__len__())
    return risultato


def creazionelista(file):
    with open(file, encoding="UTF-8", mode='r') as file:
        raw = list()
        raw += [x.split() for x in file if x[0] != "#" and len(x.strip()) > 0]
        finale = list()
        for x in raw:
            finale += [add for add in x]
    return finale


def findelemento(istruz, clu, citta):
    possibili_clues = [elemento for elemento in istruz if elemento[0] == citta and elemento[1] == clu]
    return possibili_clues


def recursive(istruz, citta, lvl, clues, frase: str, finalList, lencl):
    if lvl < lencl:
        clu = clues[lvl]
    else:
        finalList.add(tuple([frase[:-1], citta]))
        return finalList
    possibili_clues = findelemento(istruz, clu, citta)
    if possibili_clues.__len__() == 0:
        return finalList
    else:
        for indizio in possibili_clues:
            finalList = recursive(istruz, indizio[2], lvl+1, clues, frase + indizio[3] + " ", finalList, lencl)
    return finalList


def findele(ele):
    elemento = []
    ltt = str()
    maius = True
    for lettera in ele:
        stato = lettera.isupper()
        if maius:
            if stato:
                ltt += lettera
            else:
                elemento.append(ltt)
                ltt = lettera
                maius = False
        else:
            if not stato:
                ltt += lettera
            else:
                elemento.append(ltt)
                ltt = lettera
                maius = True
    elemento.append(ltt)
    return elemento


def trasforma(istruz):
    nuoval = [findele(x) for x in istruz]
    return nuoval
