package server

import "fmt"

type Vector2 struct {
	X int32
	Y int32
}

type ColorType struct {
	R int32
	G int32
	B int32
	A int32
}

const (
	TypeLength  int32 = 8
	PMessage    int8  = 0x00000000
	DMessage    int8  = 0x00000001
	ExitMessage int8  = 0x00000002
)

/* Interface to generalize messaging */
type GMessage interface {
	Output()
}

/* Point data that will be passed throught the network */
type PointMessage struct {
	Position  Vector2
	Thickness int32
	Color     ColorType
}

func (p PointMessage) Output() {
	fmt.Println(p.Position, p.Color, p.Thickness)
}
