package main

import (
	"github.com/gorilla/websocket"
)

//clientはチャットを行っている１人のユーザーを表す
type client struct {
	//socketはこのクライアントのためのWebsocket
	socket *websocket.Conn
	//sendはメッセージが送られるチャネル
	send chan []byte
	//roomはこのクライアントが参加しているチャットルーム
	room *room
}
