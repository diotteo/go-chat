package libs

import (
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
				TypeId: TypeId_REGISTER_MESSAGE_ID,
				Name: self.name,
				SentTs: timestamppb.Now(),
				}
	m := &RegisterMessage{
			Header: &h,
			}
	data, err := proto.Marshal(m)
	return data, err
}

func (self *Client) GetQuitMessage() ([]byte, error) {
	h := Header{
				TypeId: TypeId_QUIT_MESSAGE_ID,
				Name: self.name,
				SentTs: timestamppb.Now(),
				}
	m := &QuitMessage{
			Header: &h,
			}
	data, err := proto.Marshal(m)
	return data, err
}

func (self *Client) GetSendMessage(msg string) ([]byte, error) {
	h := Header{
				TypeId: TypeId_SEND_MESSAGE_ID,
				Name: self.name,
				SentTs: timestamppb.Now(),
				}
	m := &SendMessage{
			Header: &h,
			UserMessage: msg,
			}
	data, err := proto.Marshal(m)
	return data, err
}


type GenericMessage struct {
	Id TypeId
	Register *RegisterMessage
	Send *SendMessage
	Quit *QuitMessage
}

func MessageFromBytes(data []byte) (GenericMessage) {
	msg_type := TypeId(data[0])

	return GenericMessage{ Id: msg_type }
}
