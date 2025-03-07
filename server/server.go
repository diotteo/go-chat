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
			msg := libs.MessageFromBytes(buf[:n])
			switch msg.Payload.(type) {
			case *libs.GenericMessage_Register:
				client, found := clients[msg.Header.Name]
				_ = client
				if found {
					fmt.Printf("Name %s already taken\n", msg.Header.Name)
				} else {
					fmt.Printf("Registering new user %s from %s\n", msg.Header.Name, addr.IP)
					sender := libs.NewClient(msg.Header.Name, addr)
					for _, client := range clients {
						data, _ := sender.GetRegisterMessage()
						conn.WriteToUDP(data, client.Addr())
					}
					clients[msg.Header.Name] = sender
				}
			case *libs.GenericMessage_Send:
				user_msg := msg.GetSend().UserMessage
				fmt.Printf("[%s] %s\n", msg.Header.Name, user_msg)
				sender, ok := clients[msg.Header.Name]
				if ok {
					for name, client := range clients {
						data, _ := sender.GetSendMessage(user_msg)
						conn.WriteToUDP(data, client.Addr())
						_ = name
					}
				}
			case *libs.GenericMessage_Quit:
				sender, ok := clients[msg.Header.Name]
				if ok {
					fmt.Printf("** %s has left **\n", msg.Header.Name)
					delete(clients, msg.Header.Name)
					for _, client := range clients {
						data, _ := sender.GetQuitMessage()
						conn.WriteToUDP(data, client.Addr())
					}
				}
			}
		}
	}
}
