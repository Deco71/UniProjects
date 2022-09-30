# -*- coding: utf-8 -*-
def ex1(poem_filename):
    prosodia = []
    es = []
    finali = []
    contaes = 0
    roba = {}
    diz = {'e': True, 'a': True, 'i': True, 'o': True, 'n': False, 'r': False, 't': False, 'l': False, 's': False, 'u': True, 'b': False, 'c': False, 'm': False, 'p': False, 'f': False, 'h': False, 'v': False, 'w': False, 'y': True, 'd': False, 'g': False, 'q': False, 'z': False, 'j': True, 'x': False, 'k': False}
    table = {224: 97, 225: 97, 226: 97, 227: 97, 228: 97, 229: 97, 232: 101, 233: 101, 234: 101, 235: 101, 236: 105, 237: 105,
     238: 105, 239: 105, 242: 111, 243: 111, 244: 111, 245: 111, 246: 111, 248: 111, 249: 117, 250: 117, 251: 117,
     252: 117, 253: 121, 255: 121, 33: None, 39: None, 44: None, 46: None, 58: None, 59: None, 63: None, 64: None,
     35: None, 36: None, 37: None, 94: None, 38: None, 42: None, 40: None, 41: None, 95: None, 60: None, 62: None,
     45: None, 43: None, 61: None, 91: None, 93: None, 123: None, 125: None, 124: None, 96: None, 126: None, 47: None,
     49: None, 51: None, 50: None, 52: None, 53: None, 54: None, 55: None, 56: None, 57: None, 48: None, 10: None,
     32: None}
    with open(poem_filename, mode='r', encoding='UTF-8') as stringa:
        for f in stringa:
            f = f.lower().translate(table)
            punt, lunghezza = calcolatore(f, diz)
            es.append(lunghezza)
            ae = f[-int(punt):]
            finali.append(ae)
            testo = ae + str(lunghezza)
            if testo not in roba:
                roba[testo] = contaes
                contaes += 1
            prosodia.append(roba[testo])
    return prosodia, find_periodo(prosodia), es, finali


def find_periodo(prosodia):
    lunghezza = len(prosodia)
    for count in conta_divisori(lunghezza):
        conta = count + 3
        A = normalizzatore(prosodia[0:conta])
        if A == normalizzatore(prosodia[conta:conta*2]):
            if finito(A, lunghezza, prosodia, conta):
                return conta
    return len(prosodia)


def normalizzatore(lista):
    normalista = {}
    risultato = []
    contatore = 0
    for elemento in lista:
        if elemento not in normalista:
            normalista[elemento] = contatore
            contatore += 1
        risultato.append(normalista[elemento])
    return risultato


def finito(lista, lunghezza, prosodia, conta):
    for count in range(lunghezza // conta):
        verifica = normalizzatore(prosodia[conta * count:conta * (count + 1)])
        if lista != verifica:
            return False
    return True


def conta_divisori(n):
    for i in range(n // 2):
        if n % (i+3) == 0:
            yield i


def calcolatore(f, diz):
    contatore = 1
    finale = 1
    cambio = diz[f[0]]
    boolultimo = cambio
    lista = f[1:]
    for parola in lista:
        lettera = diz[parola]
        if lettera == boolultimo:
            finale += 1
        else:
            if cambio:
                contatore += 1
                cambio = False
                finale = 1
                boolultimo = lettera
            else:
                cambio = True
                finale += 1
                boolultimo = lettera
    return finale, contatore
