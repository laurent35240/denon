package main

import (
	"fmt"
	"net/http"
	"io/ioutil"
	"time"

	"github.com/gorilla/mux"
	"log"
	"encoding/json"
	"golang.org/x/net/websocket"
	"github.com/rs/cors"

	"github.com/laurent35240/denon/device"
)


var denon  = device.Denon{Host: "192.168.1.9"}

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
		denon.PowerOn()
	case powerState.State == "OFF":
		denon.PowerOff()
	}
	fmt.Fprint(w, "OK")
}

func WsServer(ws *websocket.Conn)  {
	ws.Write([]byte("Connected"))
	c:= make(chan string)
	var status, oldStatus string
	for {
		go denon.GetStatus(c)
		status = <- c
		if status != oldStatus {
			ws.Write([]byte(status))
			fmt.Printf("Status %s sent through WS\n", status)
			oldStatus = status
		}
		time.Sleep(100 * time.Millisecond)
	}
}

func main() {
	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/power", handlePower).Methods("PUT")
	router.Handle("/ws", websocket.Handler(WsServer))
	c := cors.New(cors.Options{
		AllowedMethods: []string{"GET", "POST", "PUT"},
	})
	handler := c.Handler(router)
	log.Fatal(http.ListenAndServe(":8080", handler))
}
