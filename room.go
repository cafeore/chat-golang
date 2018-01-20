package main

type room struct {
	//forwardは他のクライアントに転送するためのメッセージを保持するチャネル
	forward chan []byte
}
