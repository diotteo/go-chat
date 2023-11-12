package main

import (
	"fmt"
	"net"

	"libs"
	)

func main() {
	SERVER := ":12345"
	addr, err := net.ResolveUDPAddr("udp4", SERVER)
	if err != nil {
		panic(err)
	}

	conn, err := net.ListenUDP("udp4", addr)
	if err != nil {
		panic(err)
	}
	defer conn.Close()
	fmt.Printf("Server started on %s\n", SERVER)

	clients := make(map[string]*libs.Client)
	buf := make([]byte, 1024)
	for {
		n, addr, err := conn.ReadFromUDP(buf)
		if err != nil {
			fmt.Printf("couldn't read: %s\n", err)
		} else {
			gen_msg := libs.MessageFromBytes(buf[:n])
			switch gen_msg.MessageType() {
			case libs.RegisterMessageId:
				msg := gen_msg.(*libs.RegisterMessage)
				client, found := clients[msg.Name()]
				_ = client
				if found {
					fmt.Printf("Name %s already taken\n", msg.Name())
				} else {
					fmt.Printf("Registering new user %s from %s\n", msg.Name(), addr.IP)
					sender := libs.NewClient(msg.Name(), addr)
					for _, client := range clients {
						conn.WriteToUDP(sender.GetRegisterMessage(), client.Addr())
					}
					clients[msg.Name()] = sender
				}
			case libs.SendMessageId:
				msg := gen_msg.(*libs.SendMessage)
				fmt.Printf("[%s] %s\n", msg.Name(), msg.UserMessage())
				sender, ok := clients[msg.Name()]
				if ok {
					for name, client := range clients {
						conn.WriteToUDP(sender.GetSendMessage(msg.UserMessage()), client.Addr())
						_ = name
					}
				}
			case libs.QuitMessageId:
				msg := gen_msg.(*libs.QuitMessage)
				sender, ok := clients[msg.Name()]
				if ok {
					fmt.Printf("** %s has left **\n", msg.Name())
					delete(clients, msg.Name())
					for _, client := range clients {
						conn.WriteToUDP(sender.GetQuitMessage(), client.Addr())
					}
				}
			}
		}
	}
}
