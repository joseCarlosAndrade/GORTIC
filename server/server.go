package server

import (
	"fmt"
	"net"

)

func Print() {
	fmt.Println("server side")
}

func StartServer() {
	fmt.Println("Starting server")
	listener, err := net.Listen("tcp", PORT)
	if err != nil {
		panic(err)
	}

	clientManager := ClientManager{
		clients: make(map[*Client]bool),
		broadcast: make(chan []byte),
		register: make(chan *Client),
		unregister: make(chan *Client),
	}
	
	go clientManager.Start() // start goroutine for broadcasting, register and unregister

	// main loop to accept connections
	for {
		conn , err := listener.Accept()
		if err != nil {
			fmt.Println(err)
		}
		// creating client struct for new connection
		client := &Client{
			Socket: conn,
			Data: make(chan []byte),
		}
		// channeling client to register channel
		clientManager.register <- client
		
		// send and receive goroutines
		go clientManager.Receive(client)
		go clientManager.Send(client)
	}
}