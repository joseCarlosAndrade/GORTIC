package interfaces

import (
	"bufio"
	"fmt"
	"sync"

	rl "github.com/gen2brain/raylib-go/raylib"

	"os"
	"strings"

	"github.com/joseCarlosAndrade/GORTIC/server"
)

/* Interface handler */
type DrawingBoard struct {
	CurrentWord string
	Drawing     bool
	Canva       rl.RenderTexture2D

	ColorsAvailable []rl.Color
	ActiveColor     rl.Color

	PointBuffer []server.PointMessage
	LastPoint   []int32
	LastPointR  []int32
}

type UserInterface struct {
	Board  *DrawingBoard
	Client *server.Client

	incomingMessages chan server.GMessage     // channel to handle incoming messages
	outgoingDrawing  chan server.PointMessage // channel to handle outgoing messages (apparently its not needed???/)
	// drawingMutex sync.Mutex
}

/* Go routine to handle messaging. All client -> server socket messages should be done here */
func (i *UserInterface) HandleAssyncronousMessages() {

	// for m :=   range i.outgoingDrawing {
	// 	i.Client.SendCompleteMessage(m)
	// } // single channel
	var m sync.Mutex

	for {
		// select {
		// case m, ok := <- i.outgoingDrawing:
		// 	if !ok {
		// 		fmt.Println("Outgoing channel from client ", i.Client.Socket.LocalAddr().String(), " broke. Exiting..")
		// 		return
		// 	}
		// 	i.Client.SendCompleteMessage(m)
		// // default:
		// }

		// select {
		// case m, ok := <-i.incomingMessages:
		// 	if !ok {
		// 		fmt.Println("Message channel error on HandleAssyncronousMessages")
		// 	} else {
		// 		fmt.Println("received something")
		// 		switch m.(type) {
				
		// 	}
		// }

		if i.Board.PointBuffer != nil && len(i.Board.PointBuffer) > 0 {
			m.Lock()
			for _, p := range i.Board.PointBuffer {
				// i.outgoingDrawing <-p
				if p.Thickness == 0 {
					continue
				}
				i.Client.SendCompleteMessage(p)
				fmt.Println("Sending this point: ", p)
			}
			// if i.Board.PointBuffer != nil {
			i.Board.PointBuffer = nil
			m.Unlock()
		}
	}
}

func getInputName() string {
	reader := bufio.NewReader(os.Stdin)
	msg, er := reader.ReadString('\n')
	if er != nil {
		panic(er)
	}
	return strings.TrimRight(msg, "\n")
}

func InitInterface(drawing bool) {
	fmt.Println("\n\nHello! Please type your username and press enter (no spaces!):\n ")
	userName := getInputName()
	fmt.Println("Name chosen: ", userName)

	var cwo string
	if drawing {
		cwo = "Drawing!"
		// userName = "ClientDrawing"
	} else {
		cwo = "Guessing!"
		// userName = "ClientGuess"
	}

	board := &DrawingBoard{
		CurrentWord:     cwo,
		Drawing:         drawing,
		PointBuffer:     nil,
		ActiveColor:     rl.White,
		ColorsAvailable: AllColors,
		LastPoint:       make([]int32, 2),
		LastPointR:      make([]int32, 2),
	}

	board.LastPoint[0] = -1
	board.LastPoint[1] = -1

	board.LastPointR[0] = -1
	board.LastPointR[1] = -1

	client, err := server.StartInterfaceClient(userName)
	if err != nil {
		fmt.Println("Error creating client.")
		panic(err)
	}

	userInterface := UserInterface{
		Board:  board,
		Client: client,
		incomingMessages: make(chan server.GMessage),
		outgoingDrawing: make(chan server.PointMessage),
	}
	fmt.Println("chegou aqui!!!!!")
	// initializing go routines
	go client.Receive()
	// go userInterface.InitScreenRelated()
	go userInterface.HandleAssyncronousMessages()
	userInterface.InitScreenRelated()

	// for {
	// 	reader := bufio.NewReader(os.Stdin)
	// 	msg, _ := reader.ReadString('\n')

	// 	gm := server.PointMessage{
	// 		Position: server.Vector2{ X :10,  Y:11},
	// 		Color: server.ColorType{R: 100, G: 100, B: 100, A: 100},
	// 		Thickness: 3,
	// 	}

	// 	userInterface.Client.SendCompleteMessage(gm)

	// 	if m := strings.TrimRight(msg, "\n"); m == "exit" {
	// 		return
	// 	} else if m == "guess" {
	// 		board.CurrentWord = "Word Guessed!"
	// 		board.Drawing = false
	// 	}
	// }
}
