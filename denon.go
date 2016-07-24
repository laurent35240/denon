package main

import (
	"net"
	"fmt"
	"net/http"
	"io/ioutil"
	"encoding/xml"
	"time"

	"github.com/gorilla/mux"
	"log"
	"encoding/json"
	"golang.org/x/net/websocket"
)

type Denon struct {
	host string
	conn net.Conn
}

func (denon *Denon) connect() {
	conn, err := net.Dial("tcp", denon.host + ":23")
	if err != nil {
		fmt.Printf("Error connecting to denon system: %v", err)
	}
	denon.conn = conn
}

func (denon *Denon) sendCmd(cmd string)  {
	if (denon.conn == nil) {
		denon.connect()
	}
	cmdForConn := cmd + "\r"
	fmt.Fprint(denon.conn, cmdForConn)
	fmt.Printf("Command %s sent\n", cmd)
}

func (denon *Denon) powerOn()  {
	denon.sendCmd("PWON")
}

func (denon *Denon) powerOff()  {
	denon.sendCmd("PWSTANDBY")
}

func (denon *Denon) getStatus(c chan string)  {
	resp, err := http.Get("http://" + denon.host + "/goform/formMainZone_MainZoneXmlStatus.xml")
	if err != nil {
		fmt.Printf("Error getting status: %v", err)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("Error while reading status body: %v", err)
	}

	type Result struct {
		XMLName xml.Name `xml:"item"`
		PowerState string `xml:"Power>value"`
	}
	result := Result{PowerState: "unknown"}
	err = xml.Unmarshal(body, &result)
	if err != nil {
		fmt.Printf("Error while parsing xml: %v", err)
	}
	fmt.Printf("Power status: %s\n", result.PowerState)
	c <- result.PowerState
}

func (denon *Denon) showStatus()  {
	c:= make(chan string)
	for {
		go denon.getStatus(c)
		time.Sleep(100 * time.Millisecond)
	}
}

var denon  = Denon{host: "192.168.1.9"}

func handlePower(w http.ResponseWriter, r *http.Request)  {
	type PowerState struct {
		State string `json:"state"`
	}

	var powerState PowerState
	defer r.Body.Close()
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		panic(err)
	}
	fmt.Printf("Body: %s", body)
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	if err := json.Unmarshal(body, &powerState); err != nil {
		w.WriteHeader(422) // unprocessable entity
		if err := json.NewEncoder(w).Encode(err); err != nil {
			panic(err)
		}
	}
	fmt.Printf("Power state received: %s", powerState.State)
	switch {
	case powerState.State == "ON":
		denon.powerOn()
	case powerState.State == "OFF":
		denon.powerOff()
	}
	fmt.Fprint(w, "OK")
}

func WsServer(ws *websocket.Conn)  {
	ws.Write([]byte("Connected"))
	c:= make(chan string)
	for {
		go denon.getStatus(c)
		ws.Write([]byte( <- c))
		time.Sleep(100 * time.Millisecond)
	}
}

func main() {
	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/power", handlePower).Methods("PUT")
	router.Handle("/ws", websocket.Handler(WsServer))
	log.Fatal(http.ListenAndServe(":8080", router))
}
