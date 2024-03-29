package server

import (
	"bytes"
	"encoding/gob"
	"fmt"

)

// handle encoding on gob serialization and desserialization

func SerializeMessageData(data interface {}) ([]byte, error) { // serializing messages 
	var buffer bytes.Buffer
	encoder := gob.NewEncoder(&buffer)

	err := encoder.Encode(data)
	if err != nil {
		return nil, err
	}

	return buffer.Bytes(), nil
}

func DesserializeMessageData(data []byte, mtype []byte) (GMessage, error) { // desserializing into GMessage type
	
	reader := bytes.NewReader(data)
	decoder := gob.NewDecoder(reader)

	switch mtype[0] {

	case byte(PMessage): // desserializing to point message type (i dont know if this is ugly or not)
		var message PointMessage
		decoder.Decode(&message)
		fmt.Println("Desserializing to point message")
		return message, nil

	case byte(DMessage):
	
	case byte(RegMessage):
		var msg RegisterMessage
		decoder.Decode(&msg)
		return msg, nil

	case byte(RegSucMessage):
		var msg RegisterSuccessMessage
		decoder.Decode(&msg)
		return msg, nil

	case byte(RegFailMessage):
		var msg RegisterFailureMessage
		decoder.Decode(&msg)
		return msg, nil
	
	case byte(BeginDrawingModeT):
		var msg BeginDrawingMessage
		decoder.Decode(&msg)
		return msg, nil
	
	case byte(StopDrawingModeT):
		var msg StopDrawingMessage
		decoder.Decode(&msg)
		return msg, nil

	default:
	}

		return nil, nil
	
}

