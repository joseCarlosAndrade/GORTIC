package server

import (
	"fmt"
	"net"
	"time"
	// "bytes"
	// "encoding/gob"
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
		names: make(map[*Client]string),
	}
	
	go clientManager.Start() // start goroutine for broadcasting, register and unregister

	// main loop to accept connections
	Accepting:
	for {
		conn , err := listener.Accept()
		if err != nil {
			fmt.Println(err)
		}

		// awaiting register information
		
		// creating client struct for new connection
		client := &Client{
			Socket: conn,
			Data: make(chan []byte),

			
		}

		msg, err := ReceiveSingleMessage(client)
		if err == nil {
			// fmt.Println("[SERVER] Client accepted and read the first message")
			//register client

			switch t := msg.(type) {
			case RegisterMessage: 
				for _, m := range clientManager.names {
					if m == t.UserName { // name already used
						fmt.Println("Name '", m, "' already used. Responding request with RegisterFailureMessage. Closing this connection.")
						if e:= client.SendCompleteMessage(RegisterFailureMessage{Cause: "NAME ALREADY USED"}); e != nil {
							fmt.Println("[SERVER] Did not respond successfully: ", e )
						}
						time.Sleep(100*time.Millisecond)
						conn.Close()
						continue Accepting
					}

				}

				// new name
				fmt.Println("Successfully registered client '", t.UserName, "'")
				clientManager.names[client] = t.UserName // adding to name map
				client.SendCompleteMessage(RegisterSuccessMessage{})
			default:
				fmt.Println("[SERVER] Client ", conn.LocalAddr().String(), " did not send register information. Closing..")
				client.SendCompleteMessage(RegisterFailureMessage{Cause: "MISSING REGISTER MESSAGE"})
				conn.Close()
				continue
			}
		} else {
			// send error message and ignore client
			client.SendCompleteMessage(RegisterFailureMessage{ Cause : "ERROR ON REGISTRATION"})
			conn.Close()
			continue
		}
		// channeling client to register channel
		clientManager.register <- client
		
		// send and receive goroutines
		go clientManager.Receive(client)
		go clientManager.Send(client)
	}
}


