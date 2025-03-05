package main

import (
	"fmt"
	"net"
	"bufio"
	"os"
	str "strings"
	"strconv"
	ui "github.com/gizak/termui/v3"
	"github.com/gizak/termui/v3/widgets"
	//"flag"

	"libs"
	)

const LINE_END = '\n'

func user_input(rdr *bufio.Reader, ch chan string) {
	for {
		line, err := rdr.ReadString(byte(LINE_END))
		if err != nil {
			panic(err)
		}
		msg, _ := str.CutSuffix(line, string(LINE_END))
		ch <- msg
	}
}

func receive_messages(conn net.Conn, ch chan string, quit_ch chan bool) {
	buf := make([]byte, 1024)

	for {
		n, err := conn.Read(buf)
		if err != nil {
			fmt.Println("got an error")
			select {
			case <- quit_ch:
				fmt.Println("got quit")
				return
			default:
				fmt.Println("no quit, panicking")
				panic(err)
			}
			fmt.Println("Shouldn't run")
		}
		gen_msg := libs.MessageFromBytes(buf[:n])
		switch gen_msg.Id {
		case libs.TypeId_SEND_MESSAGE_ID:
			msg := gen_msg.Send
			ch <- fmt.Sprintf("[%s] %s", msg.Header.Name, msg.UserMessage)
		case libs.TypeId_QUIT_MESSAGE_ID:
			msg := gen_msg.Quit
			ch <- fmt.Sprintf("** %s has left **", msg.Header.Name)
		case libs.TypeId_REGISTER_MESSAGE_ID:
			msg := gen_msg.Register
			ch <- fmt.Sprintf("** %s has joined **", msg.Header.Name)
		}
	}
}

func main() {
	args := os.Args[1:]

	host := "localhost"
	port := 12345
	username := "anon"

	switch len(args) {
	case 1:
		username = args[0]
	case 2:
		username = args[0]
		host = args[1]
	case 3:
		username = args[0]
		host = args[1]
		port, _ = strconv.Atoi(args[2])
	default:
		fmt.Printf("Usage: %s {username} [host] [port]\n", os.Args[0])
		return
	}

	if err := ui.Init(); err != nil {
		panic(err)
	}
	defer ui.Close()

	chat_w, chat_h := ui.TerminalDimensions()
	chat_h -= 1

	l := widgets.NewList()
	l.Title = "Chat"
	l.SetRect(0, 0, chat_w, chat_h)
	l.Rows = make([]string, chat_h-2)
	ui.Render(l)

	p := widgets.NewParagraph()
	p.SetRect(0, chat_h, chat_w, chat_h+1)
	p.Border = false
	p.Text = fmt.Sprintf("[%s] ", username)
	ui.Render(p)

	server_url := fmt.Sprintf("%s:%d", host, port)

	user_ch := make(chan string)
	msg_ch := make(chan string)
	quit_ch := make(chan bool, 1)

	conn, err := net.Dial("udp4", server_url)
	if err != nil {
		panic(err)
	}
	defer func(conn net.Conn, quit_ch chan bool) {
		quit_ch <- true
		conn.Close()
	}(conn, quit_ch)
	//fmt.Printf("Connected to server %s\n", server_url)

	client := libs.NewClient(username, nil)
	data, _ := client.GetRegisterMessage()
	conn.Write(data)

	//rdr := bufio.NewReader(os.Stdin)
	//go user_input(rdr, user_ch)
	go receive_messages(conn, msg_ch, quit_ch)
	uiEvents := ui.PollEvents()
	msg := ""
	msg_count := 0
	b_continue := true
	for b_continue {
		select {
		case ev := <-uiEvents:
			switch ev.ID {
			case "<C-c>":
				data, _ := client.GetQuitMessage()
				conn.Write(data)
				b_continue = false
				break
			case "<Enter>":
				if len(msg) > 0 {
					data, _ := client.GetSendMessage(msg)
					conn.Write(data)
					msg = ""
				}
			case "<Backspace>":
				if len(msg) > 0 {
					msg = msg[:len(msg)-1]
				}
			case "<Space>":
				msg += " "
			default:
				if len(ev.ID) == 1 || ev.ID[0] != '<' {
					msg += ev.ID
				} else {
				}
			}
			p.Text = fmt.Sprintf("[%s] %s", username, msg)
			ui.Render(p)
		case msg := <-msg_ch:
			//fmt.Printf("Received %s\n", msg)
			if msg_count >= cap(l.Rows) {
				l.Rows = append(l.Rows, msg)
			} else {
				l.Rows[msg_count] = msg
			}
			msg_count++
			ui.Render(l)
		case msg := <-user_ch:
			data, _ := client.GetSendMessage(msg)
			conn.Write(data)
		}
	}
}
