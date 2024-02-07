package interfaces

import (
	// "fmt"

	"fmt"

	rl "github.com/gen2brain/raylib-go/raylib"
	"github.com/joseCarlosAndrade/GORTIC/server"
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
	
	userInt.Board.Canva = rl.LoadRenderTexture(ScreenWidth, ScreenHeight)
	userInt.CheckDrawing()
	fmt.Println("done")
}

func (userInt * UserInterface) CheckDrawing() { // handles the drawing part
	
	defer rl.CloseWindow()
	for !rl.WindowShouldClose() { // while is Drawing
		pos := rl.GetMousePosition()
		posX := int32(pos.X)
		posY := int32(pos.Y)

		// rl.BeginDrawing()
		// rl.BeginTextureMode(board.Canva) //  TODO: fix sudden black screen for no reason

		if userInt.Board.Drawing {
			
			if  rl.IsMouseButtonDown(rl.MouseButtonLeft) && rl.CheckCollisionPointRec(
				rl.GetMousePosition(),
				rl.NewRectangle(
					0, 0,
					float32(ScreenWidth),
					float32(ScreenHeight),
				),
			)  { // click and inside screen
				
				pointdata := server.PointMessage{
					Origin: userInt.Client.Socket.LocalAddr().String(),
					Position : server.Vector2{X: posX, Y: int32(posY)},
					Thickness: 1,
					Color: RaylibColor2SColo(userInt.Board.ActiveColor),
				}
				
				// rl.BeginTextureMode(userInt.Board.Canva)

				if userInt.Board.LastPoint[0] == -1 {
					// rl.DrawCircle(posX, posY, 1, rl.White)
					pointdata.NewLocation = true
					
				} else {
					// rl.DrawLine(userInt.Board.LastPoint[0], userInt.Board.LastPoint[1], posX, posY, rl.White)
					pointdata.NewLocation = false
				}

				userInt.Board.LastPoint[0] = posX
				userInt.Board.LastPoint[1] = posY
				
				// rl.EndTextureMode()

				// userInt.drawingMutex.Lock() // avoid racing conditions
				
				// userInt.outgoingDrawing <- pointdata
				// note: for some reason i cant just send pointdata to outgoingDrawing channel from here, it blocks 
				// the drawing even when using mutex
				userInt.Board.PointBuffer = append(userInt.Board.PointBuffer, pointdata) // points are buffered on PointBuffer to centralize socket messaging
				fmt.Println("point appendend: ", pointdata)
				// userInt.drawingMutex.Unlock()

			} else {
				userInt.Board.LastPoint[0] = -1
			}
			
		
		}
		// } else if !userInt.Board.Drawing {
		// 	fmt.Println("loop not drawing")
			
		Incoming:
			for {
				select {
				case pm := <- userInt.Client.IncomingDrawing: // TODO: change this, try a bufferized array or slice
					switch p := pm.(type) {
					case server.PointMessage:
						rl.BeginTextureMode(userInt.Board.Canva)
						// userInt.drawingMutex.Lock()

						if p.NewLocation { // new line starting at a new location
							rl.DrawCircle(
								p.Position.X,
								p.Position.Y,
								float32(p.Thickness),
								SColor2RaylibColor(p.Color),
								// rl.NewColor(uint8(p.Color.R), uint8(p.Color.G), uint8(p.Color.B), uint8(p.Color.A)),
							)

						} else { //continue previous line
							rl.DrawLine(
								userInt.Board.LastPointR[0], 
								userInt.Board.LastPointR[1],
								p.Position.X, 
								p.Position.Y, 
								SColor2RaylibColor(p.Color),
								// rl.NewColor(uint8(p.Color.R), uint8(p.Color.G), uint8(p.Color.B), uint8(p.Color.A))
							)

						}

						userInt.Board.LastPointR[0] = p.Position.X
						userInt.Board.LastPointR[1] = p.Position.Y
						// userInt.drawingMutex.Unlock()
						rl.EndTextureMode()

					case server.BeginDrawingMessage:
						userInt.Board.Drawing = true

					case server.StopDrawingMessage:
						userInt.Board.Drawing = false
						
					}
					
					// p, ok := pm.(server.PointMessage)
					
					
				default:
					break Incoming
				}
			}
			
		// } 	

		// drawing texture buffer
		rl.BeginDrawing()
		rl.ClearBackground(rl.Black)
		rl.DrawTextureRec( // i have to use this function so that the y axis is inverted
			userInt.Board.Canva.Texture, 
			rl.NewRectangle(0, 0, float32(userInt.Board.Canva.Texture.Width), -float32(userInt.Board.Canva.Texture.Height)),
			rl.NewVector2(0, 0),
			rl.White,
		)

		if userInt.Board.Drawing {
			userInt.drawColorChoices()			
		}

		// drawing mouse pointer as a reference
		rl.DrawCircle(posX, posY, 3, userInt.Board.ActiveColor)

		// draw text related
		rl.DrawText(userInt.Board.CurrentWord, 15, 30, 20, rl.Beige)
		
		rl.EndDrawing()


	}
}

func (inter *UserInterface) drawColorChoices() {
	// rectangle
	rl.DrawRectangle(ColorPositionX-15, ColorPositionY-15, int32((len(inter.Board.ColorsAvailable)-2)*30 +60), 30, rl.Beige)

	for i, r := range inter.Board.ColorsAvailable{
		var thic  float32 = 10
		if (rl.CheckCollisionPointCircle(
			rl.GetMousePosition(),
			rl.NewVector2( float32(ColorPositionX+ int32(i*30)), float32(ColorPositionY)),
			thic,
		)) {
			thic = 15

			if (rl.IsMouseButtonPressed(rl.MouseButtonLeft)) {
				inter.Board.ActiveColor = r
			}
		}

				
		rl.DrawCircle(ColorPositionX+ int32(i*30), ColorPositionY, thic, r)
	}
}


// util functions 
func SColor2RaylibColor(s server.ColorType) rl.Color {
	return rl.NewColor(uint8(s.R), uint8(s.G), uint8(s.B), uint8(s.A))
}

func RaylibColor2SColo(r rl.Color) server.ColorType {
	return server.ColorType{R :int32(r.R), G: int32(r.G), B: int32(r.B), A: int32(r.A)}
}