package server

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"io"
	"net"
)

// holding data about client and server managers

const (
	MESSAGE_LENGTH int = 4096
	PORT string = ":6700"
)

type ClientManager struct {
	clients    map[*Client]bool
	broadcast  chan []byte
	register   chan *Client
	unregister chan *Client
}

type Client struct {
	Socket net.Conn
	Data   chan []byte

	// TODO: add info
	// name string
}

func (m *ClientManager) Start() {
	for {
		select {
		case conn := <-m.register: // if channel register has new connection
			m.clients[conn] = true
			fmt.Println("added connection: ", conn.Socket.LocalAddr().String())

		case conn := <-m.unregister: // if channel unregister has new unregister request
			if _, exist := m.clients[conn]; exist {
				close(conn.Data)
				delete(m.clients, conn)
				fmt.Println("connection closed for socket: ", conn.Socket.LocalAddr().String())
			}

		case msg := <-m.broadcast: // if theres something on the broadcast channel
			for connection := range m.clients { // loop through every client
				select {
				case connection.Data <- msg: // if client is on, send message to data channel on every client
				default:
					close(connection.Data) // if not, close it
					delete(m.clients, connection)
				}
			}
		}
	}
}

func (m *ClientManager) Receive(client *Client) { // receive goroutine that exists for every active client
	for { // receives 8 bytes of information regarding the message type
		msgType := make([]byte, 1)
		_, err := client.Socket.Read(msgType)
		if err !=nil {
			fmt.Println("Error on receiving type")
			fmt.Println(err.Error())
			if err == io.EOF { // closing connection in case of a eof received
				m.unregister <- client
				break
			}
		}
		fmt.Println("Type received: ", msgType)

		msg := make([]byte, MESSAGE_LENGTH)
		length, err := client.Socket.Read(msg)
		
		if err != nil {
			fmt.Println("Error on socket ", client.Socket.LocalAddr().String(), ". Error: ")
			fmt.Println(err.Error())
		} 
		if length > 0 {
			msgdecoded,_ := DesserializeMessageData(msg, msgType)
			fmt.Println("Message received: ", msgdecoded)

			// TODO: set name operation
			// if msg == set_name client.name etc etc and not broadcast
			
			m.broadcast <- append(msgType, msg...) // weird way of appending two slices ???
			// m.broadcast <- msg // broadcasting (channeling received msg to broadcast channel)
		}
	}
}

func (m *ClientManager) Send(client *Client) { // send goroutine that send the data from client.data to the client itself
	defer client.Socket.Close() // we must use defer here so that when an error occurs, it still closes even though the function is returned
	for {
		select {
		case message, ok := <-client.Data:
			if !ok {
				// error, return and then close (defer)
				return
			}
			_, err := client.Socket.Write(message)
			if err != nil {
				fmt.Println("Could not write to ", client.Socket.LocalAddr().String())
				fmt.Println(err.Error())
				return
			}
		}
	}
}

func SerializeMessageData(data interface {}) ([]byte, error) { // serializing messages 
	var buffer bytes.Buffer
	encoder := gob.NewEncoder(&buffer)

	err := encoder.Encode(data)
	if err != nil {
		return nil, err
	}

	return buffer.Bytes(), nil
}

func DesserializeMessageData(data []byte, mtype []byte) (GMessage, error) { // desserializing into GMessage type
	
	reader := bytes.NewReader(data)
	decoder := gob.NewDecoder(reader)

	switch mtype[0] {

	case byte(PMessage): // desserializing to point message type (i dont know if this is ugly or not)
		var message PointMessage
		decoder.Decode(&message)
		fmt.Println("Desserializing to point message")
		return message, nil

	case byte(DMessage):
	default:
	}

		return nil, nil
}