package interfaces

import (
	// rl "github.com/gen2brain/raylib-go/raylib"
	"bufio"
	"fmt"
	// "fmt"
	"os"
	"strings"

	"github.com/joseCarlosAndrade/GORTIC/server"
)





func InitInterface() {

	board := &DrawingBoard{
		CurrentWord: "This word!",
		Drawing: true,
		
	}

	client, err := server.StartInterfaceClient()
	if err != nil {
		fmt.Println("Error creating client.")
		panic(err)
	}

	userInterface := UserInterface{
		Board: board,
		Client: client,
	}
	 
	// initializing go routines
	go client.Receive()
	go userInterface.InitScreenRelated()
	

	 
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
