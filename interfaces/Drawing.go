package interfaces

import rl "github.com/gen2brain/raylib-go/raylib"

const (
	ScreenWidth  int32 = 800
	ScreenHeight int32 = 600
	FPS          int32 = 120
)

/* Point data that will be passed throught the network */
type PointData struct {
	X int32
	Y int32
	Thickness int32
	Color rl.Vector3
}

/* Interface handler */
type DrawingBoard struct {
	CurrentWord string
	Drawing bool
	Canva rl.RenderTexture2D
}

func (board * DrawingBoard) InitScreenRelated() {
	board.InitScreen()
	board.CheckDrawing()
	board.Canva = rl.LoadRenderTexture(ScreenWidth, ScreenHeight)
}

func (board * DrawingBoard) InitScreen()  {
	rl.InitWindow(ScreenWidth, ScreenHeight, "Your Board!")
	rl.SetTargetFPS(FPS)

	rl.BeginTextureMode(board.Canva)
	rl.ClearBackground(rl.White)
	rl.EndTextureMode()
}

func (board * DrawingBoard) CheckDrawing() { // handles the drawing part
	defer rl.CloseWindow()
	for !rl.WindowShouldClose() { // while is Drawing

		rl.BeginDrawing()
		// rl.BeginTextureMode(board.Canva) //  TODO: fix sudden black screen for no reason

		if board.Drawing && rl.IsMouseButtonDown(rl.MouseButtonLeft) {
			pos := rl.GetMousePosition()

			rl.DrawCircle(int32(pos.X), int32(pos.Y), 4, rl.Black)
		} else if !board.Drawing {
			rl.ClearBackground(rl.White)
		} 	
		rl.DrawText(board.CurrentWord, 15, 30, 20, rl.Black)

		// rl.EndTextureMode()
		// rl.BeginDrawing()
		// rl.DrawTexture(board.Canva.Texture, 0, 0, rl.Blank) 
		rl.EndDrawing()
		// rl.BeginDrawing()
		// rl.DrawTextureRec(board.Canva.Texture, rl.Rectangle{0, 0, float32(board.Canva.Texture.Width), float32(-board.Canva.Texture.Height)}, rl.Vector2{0, 0}, rl.White)
		// rl.EndDrawing()

	}
}

