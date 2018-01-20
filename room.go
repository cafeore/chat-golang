package main

import (
	"github.com/gorilla/websocket"
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
