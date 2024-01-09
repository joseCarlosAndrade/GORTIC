package client

import (
	"bufio"
	"fmt"
	"io"
	"net"
	"os"
	"strings"

	"github.com/joseCarlosAndrade/GORTIC/server"
)

type Client server.Client // new type to implement local functions

func (client *Client)Receive() { // goroutine for client receiving
	for {
		msgType := make([]byte, 1)
		
		_, err := client.Socket.Read(msgType)
		if err !=nil {
			fmt.Println("Error on receiving type: ")
			fmt.Println(err.Error())

			if err == io.EOF {
				fmt.Println("EOF received.. closing server connection")
				break
			}
		}
		fmt.Println("Type received: ", msgType)

		msg := make([]byte, server.MESSAGE_LENGTH)
		length, err := client.Socket.Read(msg)

		if err !=nil {
			client.Socket.Close()
			fmt.Println("Error reading on socket: ", client.Socket.LocalAddr().String(), ". Closing..")
			fmt.Println(err.Error())
			panic(err)
		} 
		if length >0 {
			messageg, err  := server.DesserializeMessageData(msg, msgType)
			if err != nil {
				fmt.Println("Error desserializing message ")
			} else {
				fmt.Println("Message received: ", messageg)
			}
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

	defer client.Socket.Close()
	go client.Receive()

	for {
		reader := bufio.NewReader(os.Stdin)
		msg, _ := reader.ReadString('\n')
		strings.Split(msg, "") // let it here for now
		gm := server.PointMessage{
			Position: server.Vector2{ X :10,  Y:11},
			Color: server.ColorType{R: 100, G: 100, B: 100, A: 100},
			Thickness: 3,
		}
		serialized, err := server.SerializeMessageData(gm)
		if err != nil {
			fmt.Println("error serializing object")
		} else {
			connection.Write([]byte{byte(server.PMessage)})
			connection.Write(serialized)
		}
		// connection.Write([]byte(strings.TrimRight(msg, "\n")) )// sends msg + \n to socket
	}
}