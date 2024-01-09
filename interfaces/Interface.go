package interfaces

import (
	// rl "github.com/gen2brain/raylib-go/raylib"
	"bufio"
	// "fmt"
	"os"
	"strings"
)



func InitInterface() {
	board := &DrawingBoard{
		CurrentWord: "This word!",
		Drawing: true,
		
	}

	go board.InitScreenRelated()

	for {
		reader := bufio.NewReader(os.Stdin)
		msg, _ := reader.ReadString('\n')

		if m := strings.TrimRight(msg, "\n"); m == "exit" {
			return
		} else if m == "guess" {
			board.CurrentWord = "Word Guessed!"
			board.Drawing = false
		}
	}
}
