package main

import (
	_ "fmt"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
)

type room struct {
	//forwardは他のクライアントに転送するためのメッセージを保持するチャネル
	forward chan []byte
	//joinはチャットルームに参加しようとしているクライアントのためのチャネル
	join chan *client
	//leaveはチャットルームから体質しようとしているクライアントのためのチャネル
	leave chan *client
	//clientsには在室しているすべてのクライアントが保持される
	clients map[*client]bool
}

func (r *room) run() {
	for {
		select {
		case client := <-r.join:
			//参加
			r.clients[client] = true
			//fmt.Println("joined:", client)
		case client := <-r.leave:
			//退室
			delete(r.clients, client)
			close(client.send)
		case msg := <-r.forward:
			//すべてのクライアントにメッセージを送信
			for client := range r.clients {
				select {
				case client.send <- msg:
					//メッセージを送信
				default:
					//送信に失敗
					delete(r.clients, client)
					close(client.send)
				}
			}
		}
	}
}

func (c *client) read() {
	for {
		if _, msg, err := c.socket.ReadMessage(); err == nil {
			c.room.forward <- msg
		} else {
			break
		}
	}
	c.socket.Close()
}

func (c *client) write() {
	for msg := range c.send {
		if err := c.socket.WriteMessage(websocket.TextMessage, msg); err != nil {
			break
		}
	}
	c.socket.Close()
}

const (
	socketBufferSize  = 1024
	messageBufferSize = 256
)

var upgrader = &websocket.Upgrader{
	ReadBufferSize:  socketBufferSize,
	WriteBufferSize: socketBufferSize,
}

func (r *room) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	socket, err := upgrader.Upgrade(w, req, nil)
	if err != nil {
		log.Fatal("ServeHTTP:", err)
		return
	}
	client := &client{
		socket: socket,
		send:   make(chan []byte, messageBufferSize),
		room:   r,
	}
	r.join <- client
	defer func() { r.leave <- client }()
	go client.write()
	client.read()
}
