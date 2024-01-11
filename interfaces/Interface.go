package interfaces

import (
	"bufio"
	"fmt"
	"sync"

	rl "github.com/gen2brain/raylib-go/raylib"

	// "fmt"
	"os"
	"strings"

	"github.com/joseCarlosAndrade/GORTIC/server"
)

/* Interface handler */
type DrawingBoard struct {
	CurrentWord string
	Drawing bool
	Canva rl.RenderTexture2D
	PointBuffer []server.PointMessage
}

type UserInterface struct {
	Board *DrawingBoard
	Client *server.Client

	// incoming chan server.GMessage // channel to handle incoming messages
	outgoingDrawing chan server.PointMessage // channel to handle outgoing messages (apparently its not needed???/)
	drawingMutex sync.Mutex
}

func (i *UserInterface)HandleAssyncronousMessages() {
	for {
		// select {
		// case m, ok := <- i.outgoingDrawing:
		// 	if !ok {
		// 		fmt.Println("Outgoing channel from client ", i.Client.Socket.LocalAddr().String(), " broke. Exiting..")
		// 		return
		// 	}
		// 	i.Client.SendCompleteMessage(m)
		// default:
		// }
		
		if i.Board.PointBuffer != nil {
			for _, p := range i.Board.PointBuffer {
				// i.outgoingDrawing <-p
				i.Client.SendCompleteMessage(p)
			}
		// if i.Board.PointBuffer != nil {
			i.Board.PointBuffer = nil
		}
	}
}

func InitInterface(drawing bool) {

	board := &DrawingBoard{
		CurrentWord: "This word!",
		Drawing: drawing,
		PointBuffer: nil,
		// Canva: rl.LoadRenderTexture(ScreenWidth, ScreenHeight), // canva
	}

	client, err := server.StartInterfaceClient()
	if err != nil {
		fmt.Println("Error creating client.")
		panic(err)
	}

	userInterface := UserInterface{
		Board: board,
		Client: client,
		// incoming: make(chan server.GMessage),
		outgoingDrawing: make(chan server.PointMessage),

	}
	 
	// initializing go routines
	go client.Receive()
	go userInterface.InitScreenRelated()
	go userInterface.HandleAssyncronousMessages()

	 
	for {
		reader := bufio.NewReader(os.Stdin)
		msg, _ := reader.ReadString('\n')

		gm := server.PointMessage{
			Position: server.Vector2{ X :10,  Y:11},
			Color: server.ColorType{R: 100, G: 100, B: 100, A: 100},
			Thickness: 3,
		}

		userInterface.Client.SendCompleteMessage(gm)

		if m := strings.TrimRight(msg, "\n"); m == "exit" {
			return
		} else if m == "guess" {
			board.CurrentWord = "Word Guessed!"
			board.Drawing = false
		}
	}
}
