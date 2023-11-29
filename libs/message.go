package libs

import (
	"fmt"
	"time"
	"encoding/gob"
	"bytes"
	)

type MessageTypeId int8

const (
	RegisterMessageId MessageTypeId = iota
	SendMessageId
	QuitMessageId
)

type Message struct {
	TypeId MessageTypeId
	Name string
	Time time.Time
}

type Messager interface {
	MessageType() MessageTypeId
}

func (self Message) MessageType() MessageTypeId {
	return self.TypeId
}

func (self *Client) GetRegisterMessage() []byte {
	buf := bytes.Buffer{}
	enc := gob.NewEncoder(&buf)

	msg := RegisterMessage{
			Message: Message{
				TypeId: QuitMessageId,
				Name: self.name,
				Time: time.Now(),
				},
			}
	err := enc.Encode(msg)
	if err != nil {
		panic(err)
	}
	return buf.Bytes()
}

type RegisterMessage struct {
	Message
}


func (self *Client) GetSendMessage(user_msg string) ([]byte) {
	buf := bytes.Buffer{}
	enc := gob.NewEncoder(&buf)

	msg := SendMessage{
			Message: Message{
				TypeId: SendMessageId,
				Name: self.name,
				Time: time.Now(),
				},
			UserMessage: user_msg,
			}
	err := enc.Encode(msg)
	if err != nil {
		panic(err)
	}
	return buf.Bytes()
}

type SendMessage struct {
	Message
	UserMessage string
}


func (self *Client) GetQuitMessage() ([]byte) {
	buf := bytes.Buffer{}
	enc := gob.NewEncoder(&buf)

	msg := QuitMessage{
			Message: Message{
				TypeId: QuitMessageId,
				Name: self.name,
				Time: time.Now(),
				},
			}
	err := enc.Encode(msg)
	if err != nil {
		panic(err)
	}
	return buf.Bytes()
}

type QuitMessage struct {
	Message
}


func MessageFromBytes(data []byte) (Messager) {
	msg_type := MessageTypeId(data[0])
	dec := gob.NewDecoder(bytes.NewBuffer(data))

	switch msg_type {
	case RegisterMessageId:
		var msg RegisterMessage
		err := dec.Decode(&msg)
		if err != nil {
			panic(err)
		}
		return &msg
	case SendMessageId:
		var msg SendMessage
		err := dec.Decode(&msg)
		if err != nil {
			panic(err)
		}
		return &msg
	case QuitMessageId:
		var msg QuitMessage
		err := dec.Decode(&msg)
		if err != nil {
			panic(err)
		}
		return &msg
	default:
		panic(fmt.Sprintf("Unhandled valid message type: %d", msg_type))
	}

	return nil
}
