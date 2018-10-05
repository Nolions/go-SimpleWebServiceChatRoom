package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/gorilla/websocket"
)

type message struct {
	Action  int       `json:"action"`
	Message string    `json:"message"`
	Account string    `json:"account"`
	Time    time.Time `json:"time"`
}

var account string

const actionREADER = 0
const actionLOGIN = 1

func main() {
	// create connection
	conn, _, err := websocket.DefaultDialer.Dial("ws://127.0.0.1:8080/ws", nil)
	if err != nil {
		log.Println("dail error: ", err)
	}

	fmt.Println("connection success......")

	fmt.Println("Select a model:")
	fmt.Println("1. write mode")
	fmt.Println("2. read mode")
	fmt.Print("input your option : ")

	reader := bufio.NewReader(os.Stdin)
	input, err := reader.ReadString('\n')

	switch strings.TrimSpace(input) {
	case "1":
		write(conn)
		break
	case "2":
		go read(conn)
		for {

		}
	}
}

func registered(c *websocket.Conn) {
	fmt.Println("Client runing...")

	reader := bufio.NewReader(os.Stdin)
	fmt.Print("please input your acount name:")
	input, err := reader.ReadString('\n')
	if err != nil {
		fmt.Println("connecton error, ", err)
		return
	}

	msg := message{}
	msg.Account = strings.TrimSpace(input)

	err = c.WriteJSON(msg)
	if err != nil {
		fmt.Println("Write Messafe error, ", err)
		return
	}
}

func read(c *websocket.Conn) {
	fmt.Println("====== read mode ======")
	done := make(chan struct{})
	defer close(done)

	var msg message

	for {
		// _, msg, err := c.ReadMessage()
		err := c.ReadJSON(&msg)
		if err != nil {
			log.Println("connection error, ", err)
			return
		}

		fmt.Println(msg.Account + " say:" + msg.Message + " on " + msg.Time.Format("2006-01-02 15:04:05"))
	}
}

func write(c *websocket.Conn) {
	fmt.Println("====== write mode ======")

	defer c.Close()

	fmt.Print("your account Name : ")

	account, _ := login(c)

	readerMessage := bufio.NewReader(os.Stdin)
	for {
		fmt.Print("say >>> ")
		input, err := readerMessage.ReadString('\n')
		if err != nil {
			fmt.Print("read dtat error form bufio, ", err)
			break
		}

		// fmt.Println(account)
		// fmt.Println(strings.TrimSpace(input))

		msg := message{}
		msg.Action = actionREADER
		msg.Account = account
		msg.Message = strings.TrimSpace(input)
		msg.Time = time.Now().UTC()

		// err := c.WriteMessage(websocket.TextMessage, []byte(m))
		err = c.WriteJSON(msg)

		if err != nil {
			fmt.Println("write  message fail, ", err)
			continue
		}
	}
}

func login(conn *websocket.Conn) (string, error) {
	reader := bufio.NewReader(os.Stdin)
	account, _ := reader.ReadString('\n')
	account = strings.TrimSpace(account)

	data := message{}
	data.Action = actionLOGIN
	data.Account = account

	err := conn.WriteJSON(data)
	err = conn.ReadJSON(&data)

	return data.Account, err
}
