package server

import (
	"bufio"
	"fmt"
	"io"
	"net"
	"os"
	"strings"

	// "github.com/joseCarlosAndrade/GORTIC/server"
	// "github.com/joseCarlosAndrade/GORTIC/server"
)



type Client struct {
	Socket net.Conn
	Data   chan []byte // channel for server -> client messages broadcasting

	IncomingDrawing chan GMessage
	// TODO: add info
	// name string
}

// type Client server.Client // new type to implement local functions

func (client *Client) Receive() { // goroutine for client receiving
	for {
		msgType := make([]byte, 1)

		_, err := client.Socket.Read(msgType)
		if err != nil {
			fmt.Println("Error on receiving type: ")
			fmt.Println(err.Error())

			if err == io.EOF {
				fmt.Println("EOF received.. closing server connection")
				break
			}
		}
		fmt.Println("Type received: ", msgType)

		msg := make([]byte, MESSAGE_LENGTH)
		length, err := client.Socket.Read(msg)
		fmt.Println("i am receiving!")

		if err != nil {
			client.Socket.Close()
			fmt.Println("Error reading on socket: ", client.Socket.LocalAddr().String(), ". Closing..")
			fmt.Println(err.Error())
			panic(err)
		}
		if length > 0 {
			messageg, err := DesserializeMessageData(msg, msgType)
			
			if err != nil {
				panic(err)
			}
			switch m := messageg.(type) {
			case PointMessage:
				
				// if m.Origin == client.Socket.LocalAddr().String() { // if comes from itself
				// 	continue
				// }

				client.IncomingDrawing <- m

			case BeginDrawingMessage:
				
			default:
			}
			
			if err != nil {
				fmt.Println("Error desserializing message ")
			} else {
				fmt.Println("Message received: ", messageg)
			}
		}
	}
}

/* Serialization and socket writing wrapper function implemented on *Client type */
func (client *Client) SendCompleteMessage(msg GMessage) error {
	switch msg.(type) {
	case PointMessage: // checks message type to send typing information
		client.Socket.Write([]byte{byte(PMessage)})
	case RegisterMessage:
		client.Socket.Write([]byte{byte(RegMessage)})

	case RegisterFailureMessage:
		client.Socket.Write([]byte{byte(RegFailMessage)}) 
	case RegisterSuccessMessage:
		client.Socket.Write([]byte{byte(RegSucMessage)})
		
	default:
		fmt.Println("Typing incorrect on sendcomplete message.")
		return &RegisterError{} //TODO CREATE A PROPER ERROR HERE
	}
	// connection.Write([]byte{byte(msgType)}) // sending type information

	serialized, err := SerializeMessageData(msg)
	if err != nil {
		return err
	}

	client.Socket.Write(serialized)
	return nil
}

func StartDebugClient() { // main client starter
	fmt.Println("Starting client..")
	addr := fmt.Sprintf("localhost%s", PORT)
	connection, err := net.Dial("tcp", addr)

	if err != nil {
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

		gm := PointMessage{
			Position:  Vector2{X: 10, Y: 11},
			Color:     ColorType{R: 100, G: 100, B: 100, A: 100},
			Thickness: 3,
		}
		err := client.SendCompleteMessage(gm)
		if err != nil {
			fmt.Println("error sending complete message")
		}

	}
}

func StartInterfaceClient(name string) (*Client, error) {
	fmt.Println("Dialing ...")
	addr := fmt.Sprintf("localhost%s", PORT)
	conn, err := net.Dial("tcp", addr)

	if err != nil {
		return nil, err
	}

	fmt.Println("Registering..")
	// send register information
	r := RegisterMessage{
		Origin: conn.LocalAddr().String(),
		UserName: name,
	}

	client := &Client{
		Socket:          conn,
		IncomingDrawing: make(chan GMessage, 100),
	}

	if e := client.SendCompleteMessage(r); e != nil { // registration
		return nil, e // coulnd sent
	} else {
		// awaiting server response
		fmt.Println("Sent registration protocol. Waiting for server aproval..")
		response, err := ReceiveSingleMessage(client)
		if err != nil {
			return nil, err
		}

		switch r := response.(type) {
		case RegisterSuccessMessage:
			fmt.Println("[CLIENT] Registration accepted")
			return client, nil
		case RegisterFailureMessage:
			fmt.Println("Registration failed: ", r.Cause)
			return nil, &RegisterError{}
		default: // received any message other than register
			fmt.Println("ops!")
			return nil, &RegisterError{}
		}

		
	}
	// go client.Receive()

	// return client, nil
}

func ReceiveSingleMessage(client *Client) (GMessage, error) {
	msgType := make([]byte, 1)
		_, err := client.Socket.Read(msgType)
		if err !=nil {
			fmt.Println("Error on receiving type")
			fmt.Println(err.Error())
			if err == io.EOF { // closing connection in case of a eof received
				return nil, err
				
			}
		}
		fmt.Println("Type received: ", msgType)

		msg := make([]byte, MESSAGE_LENGTH)
		length, err := client.Socket.Read(msg)
		
		if err != nil {
			fmt.Println("Error on socket ", client.Socket.LocalAddr().String(), ". Error: ")
			fmt.Println(err.Error())
			return nil, err
		} 
		if length > 0 {
			msgdecoded,_ := DesserializeMessageData(msg, msgType)
			fmt.Println(" Message received: ", msgdecoded)
			return msgdecoded, nil
			// TODO: set name operation
			// if msg == set_name client.name etc etc and not broadcast
			
			// m.broadcast <- append(msgType, msg...) // weird way of appending two slices ???
			// m.broadcast <- msg // broadcasting (channeling received msg to broadcast channel)
		} else {
			return nil, &EmptyMessageError{}
		}
}