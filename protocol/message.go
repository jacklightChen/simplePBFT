package protocol

import (
	"encoding/json"
	pb "github.com/simplePBFT/protos"
)

/*
func (msg []byte) 通过type解析message
func (msg []byte) is callback function used to unserialize message according to type (byte)
*/
//var SerialRegister  = map[byte] func (msg []byte) (Message,error){}


//type Message interface {
//	Type() byte
//	CommitteeID() string
//	SetCommitteeID(string)
//}

//func RegistUnserialHandler(t byte, unpack func (msg []byte) (Message, error)) (error) {
//	if _, err := SerialRegister[t]; err {
//		return errors.New("this message type is registed")
//	}
//	SerialRegister[t] = unpack
//	return nil
//}

func Serial(msg pb.Message) ([]byte, error) {
	jsonBytes, err := json.Marshal(msg)
	if err != nil {
		return nil, err
	}
	return jsonBytes, nil
}

func Unserial(msgBytes []byte) (pb.Message, error) {
	//preType := msgBytes[0]
	//if unpack, ok := SerialRegister[preType]; ok {
	//	return unpack(msgBytes[1:])
	//}
	//return nil, errors.New(fmt.Sprint("no this message type %x\n", preType))
	var msg pb.Message
	err := json.Unmarshal(msgBytes, &msg)
	return msg, err
}