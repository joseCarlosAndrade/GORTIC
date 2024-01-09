package interfaces

import (
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

/* Interface handler */
type DrawingBoard struct {
	CurrentWord string
	Drawing bool
	Canva rl.RenderTexture2D
}

type UserInterface struct {
	Board *DrawingBoard
	Client *server.Client
}


func (board * DrawingBoard) InitScreen()  {
	rl.InitWindow(ScreenWidth, ScreenHeight, "Your Board!")
	rl.SetTargetFPS(FPS)

	rl.BeginTextureMode(board.Canva)
	rl.ClearBackground(rl.White)
	rl.EndTextureMode()
}

func (userInt * UserInterface) InitScreenRelated() {
	userInt.Board.InitScreen()
	userInt.CheckDrawing()
	userInt.Board.Canva = rl.LoadRenderTexture(ScreenWidth, ScreenHeight)
}

func (userInt * UserInterface) CheckDrawing() { // handles the drawing part
	defer rl.CloseWindow()
	for !rl.WindowShouldClose() { // while is Drawing

		rl.BeginDrawing()
		// rl.BeginTextureMode(board.Canva) //  TODO: fix sudden black screen for no reason

		if userInt.Board.Drawing && rl.IsMouseButtonDown(rl.MouseButtonLeft)  {
			pos := rl.GetMousePosition()
			pointdata := server.PointMessage{
				Position : server.Vector2{X: int32(pos.X), Y: int32(pos.Y)},
				Thickness: 4,
				Color: server.ColorType{R: 255, G: 255, B: 255, A: 255},
			}

			rl.DrawCircle(int32(pos.X), int32(pos.Y), 4, rl.Black)
			
			userInt.Client.SendCompleteMessage(pointdata)

		} else if !userInt.Board.Drawing {
			rl.ClearBackground(rl.White)
		} 	
		rl.DrawText(userInt.Board.CurrentWord, 15, 30, 20, rl.Black)

		// rl.EndTextureMode()
		// rl.BeginDrawing()
		// rl.DrawTexture(board.Canva.Texture, 0, 0, rl.Blank) 
		rl.EndDrawing()
		// rl.BeginDrawing()
		// rl.DrawTextureRec(board.Canva.Texture, rl.Rectangle{0, 0, float32(board.Canva.Texture.Width), float32(-board.Canva.Texture.Height)}, rl.Vector2{0, 0}, rl.White)
		// rl.EndDrawing()

	}
}

