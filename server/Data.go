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
	// message typing
	TypeLength     int32 = 8
	PMessage       int8  = 0x00000000
	DMessage       int8  = 0x00000001
	ExitMessage    int8  = 0x00000002
	RegMessage     int8  = 0x00000003
	RegFailMessage int8  = 0x00000004
	RegSucMessage  int8  = 0x00000005

	// drawing commands
	EraseAll  int8 = 0x00000001
	FillColor int8 = 0x00000001
)

/* Interface to generalize messaging */
type GMessage interface {
	Output()
}

/* Point data that will be passed throught the network */
type PointMessage struct {
	Origin      string
	Position    Vector2
	Thickness   int32
	Color       ColorType
	NewLocation bool
}

func (p PointMessage) Output() {
	fmt.Println(p.Position, p.Color, p.Thickness)
}

type DrawMessage struct {
	Origin  string
	Message string
}

func (d DrawMessage) Output() {
	fmt.Println(d.Origin)
}

type DrawCommand struct {
	Origin  string
	Command int8
	Info    string
}

type RegisterMessage struct {
	Origin   string
	UserName string
}

func (r RegisterMessage) Output() {
	fmt.Println(r.Origin, r.UserName)
}

type RegisterSuccessMessage struct {
}

func (r RegisterSuccessMessage) Output() {
	fmt.Println("")
}

type RegisterFailureMessage struct {
	Cause string
}

func (r RegisterFailureMessage) Output() {
	fmt.Println(r.Cause)
}

type EmptyMessageError struct{}

func (e *EmptyMessageError) Error() string {
	return "Empty message"
}

type RegisterError struct{}

func (e *RegisterError) Error() string {
	return "Register error"
}
