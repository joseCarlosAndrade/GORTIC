package interfaces

import (
	rl "github.com/gen2brain/raylib-go/raylib"
)

const (
	ScreenWidth  int32 = 800
	ScreenHeight int32 = 600
	FPS          int32 = 120
	

	ColorPositionX int32 = 30
	ColorPositionY int32 = ScreenHeight - 30
)

var (
	AllColors = []rl.Color {
		rl.White, rl.Gray, rl.Blue, rl.Red, rl.Green, rl.Yellow,
	}
)