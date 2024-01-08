package main

import (
	"github.com/joseCarlosAndrade/GORTIC/client"
	"github.com/joseCarlosAndrade/GORTIC/server"
	"flag"
	"strings"
)

func main() {
	// client.StartClient()
	// server.StartServer()
	// help: https://www.thepolyglotdeveloper.com/2017/05/network-sockets-with-the-go-programming-language/
	flagMode := flag.String("mode", "server", "start in client or server mode")
    flag.Parse()
    if strings.ToLower(*flagMode) == "server" {
        server.StartServer()
    } else {
        client.StartClient()
    }
}	