package device

import (
	"net/http"
	"fmt"
	"io/ioutil"
	"encoding/xml"
	"time"
	"net"
)

type Denon struct {
	Host string
	conn net.Conn
}

func (denon *Denon) connect() {
	conn, err := net.Dial("tcp", denon.Host + ":23")
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

func (denon *Denon) PowerOn()  {
	denon.sendCmd("PWON")
}

func (denon *Denon) PowerOff()  {
	denon.sendCmd("PWSTANDBY")
}

func (denon *Denon) GetStatus(c chan string)  {
	resp, err := http.Get("http://" + denon.Host + "/goform/formMainZone_MainZoneXmlStatus.xml")
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
		go denon.GetStatus(c)
		time.Sleep(100 * time.Millisecond)
	}
}
