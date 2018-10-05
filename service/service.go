package main

import (
	"html/template"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
	"github.com/satori/go.uuid"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

type client struct {
	ID   string
	CONN *websocket.Conn
}

type message struct {
	Action  int       `json:"action"`
	Message string    `json:"message"`
	Account string    `json:"account"`
	Time    time.Time `json:"time"`
}

var clients []client
var users []string

func main() {
	log.Println("Service stating...")

	http.HandleFunc("/ws", wsHandleFunc)
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		httpTemplate := template.Must(template.ParseFiles("index.html"))
		httpTemplate.Execute(w, "ws://"+r.Host+"/ws")
	})

	http.ListenAndServe(":8080", nil)
}

func wsHandleFunc(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("upgrader error, ", err)
		return
	}

	//
	// creater connection
	//
	client := client{
		ID:   uuid.Must(uuid.NewV4()).String(), // get uuid
		CONN: conn,
	}
	addClient(client)

	log.Println("connection success, id:", client.ID)

	go run(client)
}

func run(c client) {
	defer c.CONN.Close()

	var msg message

	for {
		// _, m, _ := c.CONN.ReadMessage()

		// log.Println("recv: type => ", reflect.TypeOf(string(m)))
		// log.Println("recv: msg => ", string(m))

		err := c.CONN.ReadJSON(&msg)

		if err != nil {
			log.Println("read error, ", c.ID)
			log.Println("error:", err)
			removeClient(c)

			return
		}
		log.Println(msg)

		switch msg.Action {
		case 1:
			account := addUser(msg.Account)
			msg.Account = account
			c.CONN.WriteJSON(msg.Account)
			break
		case 0:
			publish(msg)
			break
		}
	}
}

// push message to all client
func publish(msg message) {
	for _, c := range clients {
		// c.CONN.WriteMessage(1, []byte(msg.Message))
		c.CONN.WriteJSON(msg)
	}
}

func addUser(account string) string {
	for _, user := range users {
		if user == account {
			log.Println(account + " come back")
			return account
		}
	}

	log.Println(account + " add the room")
	users = append(users, account)

	return account
}

// add client
func addClient(c client) {
	clients = append(clients, c)
}

// remove clinet
func removeClient(client client) {
	for i, c := range clients {
		if c.ID == client.ID {
			if len(clients) > 1 {
				clients = append(clients[:i], clients[i+1])
			} else {
				clients = clients[:cap(clients)]
			}
		}
	}
}
