package interfaces

import (
	// "fmt"

	"fmt"

	rl "github.com/gen2brain/raylib-go/raylib"
	"github.com/joseCarlosAndrade/GORTIC/server"
)

const (
	ScreenWidth  int32 = 800
	ScreenHeight int32 = 600
	FPS          int32 = 120
)

/* Point data that will be passed throught the network */
// type PointData struct {
// 	X int32
// 	Y int32
// 	Thickness int32
// 	Color rl.Vector3
// }


func (board * DrawingBoard) InitScreen()  {
	rl.InitWindow(ScreenWidth, ScreenHeight, "Your Board!")
	rl.SetTargetFPS(FPS)

	rl.BeginTextureMode(board.Canva)
	rl.ClearBackground(rl.Black)
	rl.EndTextureMode()
}

func (userInt * UserInterface) InitScreenRelated() {
	userInt.Board.InitScreen()
	userInt.CheckDrawing()
	userInt.Board.Canva = rl.LoadRenderTexture(ScreenWidth, ScreenHeight)
	fmt.Println("done")
}

func (userInt * UserInterface) CheckDrawing() { // handles the drawing part
	defer rl.CloseWindow()
	for !rl.WindowShouldClose() { // while is Drawing

		rl.BeginDrawing()
		// rl.BeginTextureMode(board.Canva) //  TODO: fix sudden black screen for no reason
		if userInt.Board.Drawing {
			if  rl.IsMouseButtonDown(rl.MouseButtonLeft)  { // click
				fmt.Println("clicking")
				pos := rl.GetMousePosition()
				pointdata := server.PointMessage{
					Origin: userInt.Client.Socket.LocalAddr().String(),
					Position : server.Vector2{X: int32(pos.X), Y: int32(pos.Y)},
					Thickness: 4,
					Color: server.ColorType{R: 10, G: 10, B: 10, A: 255},
				}
	
				rl.DrawCircle(int32(pos.X), int32(pos.Y), 4, rl.White)
				userInt.outgoingDrawing <- pointdata
				rl.EndDrawing()
				rl.BeginDrawing()
				rl.DrawCircle(int32(pos.X), int32(pos.Y), 4, rl.White)


			}
		 // sending information to server
			// fmt.Println("AAAAAAAAAAAAAAAAAAAAAAAAAA")

			// userInt.Client.SendCompleteMessage(pointdata) // send to server

		} else if !userInt.Board.Drawing {
			// for pm := range userInt.Client.IncomingDrawing {
			// 	p, ok := pm.(server.PointMessage)
			// 	if ok {
			// 		fmt.Println("drawing incoming..")
			// 		rl.DrawCircle(
			// 		p.Position.X,
			// 		p.Position.Y,
			// 		float32(p.Thickness),
			// 		rl.Blue,
			// 	)
			// 	}
			// }
			// rl.ClearBackground(rl.White)
		Incoming:
			for {
				select {
				case pm := <- userInt.Client.IncomingDrawing:
					p, ok := pm.(server.PointMessage)
					if ok {
						fmt.Println("drawing incoming..")
						rl.DrawCircle(
							p.Position.X,
							p.Position.Y,
							float32(p.Thickness),
							rl.Blue,
							)
					rl.EndDrawing()
					rl.BeginDrawing()
					rl.DrawCircle(
						p.Position.X,
						p.Position.Y,
						float32(p.Thickness),
						rl.Blue,
						)
					}
					
				default:
					break Incoming
				}
			}
			
		} 	
		
		rl.EndDrawing()
		rl.BeginDrawing()
		rl.DrawText(userInt.Board.CurrentWord, 15, 30, 20, rl.Beige)

		// rl.EndTextureMode()
		// rl.BeginDrawing()
		// rl.DrawTexture(board.Canva.Texture, 0, 0, rl.Blank) 
		rl.EndDrawing()
		// rl.BeginDrawing()
		// rl.DrawTextureRec(board.Canva.Texture, rl.Rectangle{0, 0, float32(board.Canva.Texture.Width), float32(-board.Canva.Texture.Height)}, rl.Vector2{0, 0}, rl.White)
		// rl.EndDrawing()

	}
}

