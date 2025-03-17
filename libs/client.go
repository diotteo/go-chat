package libs

import (
	"fmt"
	"net"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/timestamppb"
	)

type Client struct {
	name string
	addr *net.UDPAddr
}

func NewClient(name string, addr *net.UDPAddr) (*Client) {
	client := Client{
			name: name,
			addr: addr,
			}
	return &client
}

func (self *Client) Addr() (*net.UDPAddr) {
	return self.addr
}

func (self *Client) GetRegisterMessage() ([]byte, error) {
	h := Header{
				Name: self.name,
				SentTs: timestamppb.Now(),
				}
	m := &GenericMessage{
			Header: &h,
			Payload: &GenericMessage_Register{
				Register: &RegisterMessage{},
				},
			}
	data, err := proto.Marshal(m)
	if err != nil {
		panic(fmt.Sprintf("Failed to marshall register: %s", err))
	}
	return data, err
}

func (self *Client) GetQuitMessage() ([]byte, error) {
	h := Header{
				Name: self.name,
				SentTs: timestamppb.Now(),
				}
	m := &GenericMessage{
			Header: &h,
			Payload: &GenericMessage_Quit{
				Quit: &QuitMessage{},
				},
			}
	data, err := proto.Marshal(m)
	if err != nil {
		panic(fmt.Sprintf("Failed to marshall quit: %s", err))
	}
	return data, err
}

func (self *Client) GetSendMessage(msg string) ([]byte, error) {
	h := Header{
				Name: self.name,
				SentTs: timestamppb.Now(),
				}
	m := &GenericMessage{
			Header: &h,
			Payload: &GenericMessage_Send{
				Send: &SendMessage{
					UserMessage: msg,
					},
				},
			}
	data, err := proto.Marshal(m)
	if err != nil {
		panic(fmt.Sprintf("Failed to marshall send: %s", err))
	}
	return data, err
}


func MessageFromBytes(data []byte) (*GenericMessage) {
	msg := &GenericMessage{}
	err := proto.Unmarshal(data, msg)
	if err != nil {
		panic(fmt.Sprintf("Failed to unmarshall register: %s", err))
	}
	return msg
}
