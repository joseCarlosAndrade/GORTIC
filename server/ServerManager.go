package server

import (
	// "bytes"
	// "encoding/gob"
	"fmt"
	"io"
	"net"
	"time"
	// "net"
)

// holding data about client and server managers

const (
	MESSAGE_LENGTH int    = 4096
	PORT           string = ":6700"
)

type ClientManager struct {
	clients    map[*Client]bool
	broadcast  chan []byte
	register   chan *Client
	unregister chan *Client

	names       map[*Client]string
	playersTurn map[*Client]bool

	serverSideMessaging chan GMessage
}

/* Goroutine to handle register, unregister and broadcast channels of client manager */
func (m *ClientManager) Start() {
	for {
		select {
		case conn := <-m.register: // if channel register has new connection
			m.clients[conn] = true
			m.playersTurn[conn] = false
			// fmt.Println("added connection: ", conn.Socket.LocalAddr().String())

		case conn := <-m.unregister: // if channel unregister has new unregister request
			if _, exist := m.clients[conn]; exist {
				close(conn.Data)
				delete(m.clients, conn)
				delete(m.names, conn)
				delete(m.playersTurn, conn)
				fmt.Println("[SERVER MANAGER] Connection closed for socket: ", conn.Socket.LocalAddr().String())
				m.serverSideMessaging <- ExitMessage{name: conn.Socket.LocalAddr().String()}
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

/* Receive goroutine that runs for every online client */
func (m *ClientManager) Receive(client *Client) {
	for { // receives 8 bytes of information regarding the message type
		msgType := make([]byte, 1)
		_, err := client.Socket.Read(msgType)
		if err != nil {
			fmt.Println("[SERVER MANAGER] Error on receiving type")
			fmt.Println(err.Error())
			if err == io.EOF { // closing connection in case of a eof received
				fmt.Println("[SERVER MANAGER] EOF received.. closing")
				m.unregister <- client
				break
			} else if err == net.ErrClosed {

				m.unregister <- client
				fmt.Println("Connection broken, finishing client")
				break

			}

			if opErr, ok := err.(*net.OpError); ok && opErr.Err.Error() == "use of closed network connection" {
				fmt.Println("[SERVER MANAGER] Attempting to message a closed connection. Closing.. ")
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
			if err == net.ErrClosed {
				m.unregister <- client
				fmt.Println("Connection broke, finishing client")
				break
			}
		}
		if length > 0 {
			msgdecoded, _ := DesserializeMessageData(msg, msgType)

			switch msgType[0] {
			case byte(PMessage):
				fmt.Println("[BROADCASTING] Message received: ", msgdecoded)

				// TODO: set name operation
				// if msg == set_name client.name etc etc and not broadcast

				m.broadcast <- append(msgType, msg...) // weird way of appending two slices ???
				// m.broadcast <- msg // broadcasting (channeling received msg to broadcast channel)

			case byte(DMessage):

			case byte(RegMessage):

			case byte(RegSucMessage):
				m.serverSideMessaging <- msgdecoded

			case byte(RegFailMessage):

			case byte(BeginDrawingModeT):

			case byte(StopDrawingModeT):

			}

		}
	}
}

/* Send goroutine that send the data from client.Data channel to the client socket itself */
func (m *ClientManager) Send(client *Client) {
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

/* handles all server logic such as word choosing, current player, correct and incorrect guesses, etc*/
func (cm *ClientManager) ServerSideLogic() {
	// gameState := "halt" // halt drawing

	for {
		playersAmount := len(cm.clients)
		if playersAmount < 2 {
			// gameState = "halt"
			// fmt.Println("Waiting for players.. (min 2)")
			continue
		}
		// gameState = "drawing"
		// enough players. Find a player that hasnt played yet in this round and send a word for them to draw, while sending the others
		// to enter guess mode

		if drawer := findRemainingPlayers(cm.playersTurn); drawer != nil {
			// send word to this player and make it enter drawing mode while guessing for others
			drawer.SendCompleteMessage(BeginDrawingMessage{})

			// send the word to be drawn

			// clear everyones screen and start guessing mode
			for player := range cm.clients {
				if player == drawer {
					continue
				}
				player.SendCompleteMessage(StopDrawingMessage{})
			}

			time.Sleep(10 * time.Second) // sleeps for 10 seconds

		} else {
			// reset everyone
			resetPlayersTurns(cm.playersTurn)
		}
	}

}

// receives all clients as parameters and returns either nil if there's no player left or next drawer
func findRemainingPlayers(turns map[*Client]bool) *Client {
	for client := range turns {
		if turns[client] == false {
			//choose this as next drawer
			turns[client] = true
			return client
		}
	}
	// everyone has played already
	return nil
}

// sets all players turn to false, finishing a round
func resetPlayersTurns(turns map[*Client]bool) {
	for client := range turns {
		turns[client] = false
	}
}
