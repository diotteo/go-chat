package libs

import (
	"net"
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
