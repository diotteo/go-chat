package libs

import (
	"fmt"
	"time"
	"encoding/gob"
	)

type MessageTypeId int8

const (
	RegisterMessageId MessageTypeId = iota
	SendMessageId
	QuitMessageId
)

type Message struct {
	type_id MessageTypeId
	name string
	time time.Time
}

type Messager interface {
	MessageType() MessageTypeId
	Time() time.Time
	Name() string
}

func (c *Client) GetRegisterMessage() []byte {
	buf := bytes.Buffer{}
	raw := []byte{byte(RegisterMessageId)}
	raw = append(raw, []byte(time.Now()))
	raw = append(raw, []byte(c.name + "\x00")...)
	return raw
}

func (c *Client) GetSendMessage(msg string) ([]byte) {
	raw := append([]byte{byte(SendMessageId)}, []byte(c.name + "\x00" + msg + "\x00")...)
	return raw
}

func (c *Client) GetQuitMessage() ([]byte) {
	raw := append([]byte{byte(QuitMessageId)}, []byte(c.name + "\x00")...)
	return raw
}


type RegisterMessage struct {
	Message
}

func (self *RegisterMessage) MessageType() MessageTypeId {
	return self.type_id
}

func (self *RegisterMessage) Name() string {
	return self.name
}

func (self *RegisterMessage) Time() time.Time {
	return self.time
}


type SendMessage struct {
	Message
	message string
}

func (self *SendMessage) Name() string {
	return self.name
}

func (self *SendMessage) MessageType() MessageTypeId {
	return self.type_id
}

func (self *SendMessage) Time() time.Time {
	return self.time
}

func (self *SendMessage) UserMessage() string {
	return self.message
}


type QuitMessage struct {
	Message
}

func (self *QuitMessage) MessageType() MessageTypeId {
	return self.type_id
}

func (self *QuitMessage) Name() string {
	return self.name
}

func (self *QuitMessage) Time() time.Time {
	return self.time
}

func MessageFromBytes(data []byte) (Messager) {
	msg_type := MessageTypeId(data[0])

	STR_END := []byte("\x00")[0]

	switch msg_type {
	case RegisterMessageId:
		field_start := 1
		name := ""
		for i, b := range data[1:] {
			if b == STR_END {
				name = string(data[field_start:i+1])

				msg := RegisterMessage{
						Message: Message{
							type_id: RegisterMessageId,
							name: name,
						},
						}
				return &msg
			}
		}
		panic(fmt.Sprintf("Register message too long: %d", len(data)))
	case SendMessageId:
		name := ""
		msg_text := ""
		field_idx := 0
		field_start := 1

		for i, b := range data[1:] {
			if b == STR_END {
				switch field_idx {
				case 0:
					name = string(data[field_start:i+1])
					field_start = 1+i+1
					field_idx++
				case 1:
					msg_text = string(data[field_start:i+1])
					field_idx++

					msg := SendMessage{
							Message: Message{
								type_id: SendMessageId,
								name: name,
							},
							message: msg_text,
							}
					//fmt.Printf("Received message [%s]: %s\n", name, msg_text)
					return &msg
				}
			}
		}
		panic("SendMessageId bug encountered, please report")
	case QuitMessageId:
		name := ""
		field_idx := 0
		field_start := 1

		for i, b := range data[1:] {
			if b == STR_END {
				switch field_idx {
				case 0:
					name = string(data[field_start:i+1])
					field_start = 1+i+1
					field_idx++

					msg := QuitMessage{
							Message: Message{
								type_id: QuitMessageId,
								name: name,
							},
							}
					return &msg
				}
			}
		}
		panic("QuitMessageId bug encountered, please report")
	default:
		panic(fmt.Sprintf("Unhandled valid message type: %d", msg_type))
	}

	return nil
}
