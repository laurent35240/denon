# Denon go server

This project was kind of an exercise for controlling a denon hi-fi system with a server written in go.
It was tested with a [denon ceol N8](http://www.denon.fr/fr/product/compactsystems/networkmusicsystems/ceoln8)

## Features

* Possibility to turn it off and on through API
* Getting it power status through websocket

## Installation
* Change in `denon.go` file the IP with the local IP of your denon system
* Build the bin file `go build`
* Run the server `./denon`
* You can use the client written with [react](https://facebook.github.io/react/) which is located in `client` directory