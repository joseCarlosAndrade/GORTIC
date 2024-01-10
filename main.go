package main

import (
	// "github.com/joseCarlosAndrade/GORTIC/client"
	"github.com/joseCarlosAndrade/GORTIC/server"
	"github.com/joseCarlosAndrade/GORTIC/interfaces"
	"flag"
	"strings"
)

func main() {
	// client.StartClient()
	// server.StartServer()
	// help: https://www.thepolyglotdeveloper.com/2017/05/network-sockets-with-the-go-programming-language/
	flagMode := flag.String("mode", "server", "start in client or server mode")
    flag.Parse()
    if s := strings.ToLower(*flagMode); s == "server" {
        server.StartServer()
    } else if s == "client" {
        server.StartDebugClient()
    } else if s == "drawing" {
		interfaces.InitInterface(true)
	} else {
		interfaces.InitInterface(false)
	}
}	