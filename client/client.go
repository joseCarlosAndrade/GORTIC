package client

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strings"

	"github.com/joseCarlosAndrade/GORTIC/server"
)

type Client server.Client // new type to implement local functions

func (client *Client)Receive() { // goroutine for client receiving
	for {
		msg := make([]byte, server.MESSAGE_LENGTH)
		length, err := client.Socket.Read(msg)

		if err !=nil {
			client.Socket.Close()
			fmt.Println("Error reading on socket: ", client.Socket.LocalAddr().String(), ". Closing..")
		} 
		if length >0 {
			fmt.Println("Message received: ", string(msg))
		}
	}	
}

func StartClient() { // main client starter
	fmt.Println("Starting client..")
	addr := fmt.Sprintf("localhost%s", server.PORT)
	connection, err := net.Dial("tcp", addr)

	if err !=nil {
		panic(err)
	}

	client := &Client{
		Socket: connection,
	}

	go client.Receive()

	for {
		reader := bufio.NewReader(os.Stdin)
		msg, _ := reader.ReadString('\n')
		connection.Write([]byte(strings.TrimRight(msg, "\n")) )// sends msg + \n to socket
	}
}