package main

import (
	"time"
)

//messageは１つのメッセージ
type message struct {
	Name    string
	Message string
	When    time.Time
}
