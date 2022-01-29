package main

import (
	"fmt"
	"net"
)

type User struct {
	Name string
	Addr string

	ReceiveChan chan string
	Conn        net.Conn
}

func InitUser(conn net.Conn) *User {
	addr := conn.RemoteAddr().String()
	user := &User{
		Name:        addr + "_User",
		Addr:        addr,
		ReceiveChan: make(chan string),
		Conn:        conn,
	}
	go user.ListenBroadcastMessage()
	return user
}

func (user *User) ListenBroadcastMessage() {
	for {
		msg := <-user.ReceiveChan
		fmt.Printf("用户%s收到来自广播的消息：%s \n", user.Name, msg)
	}
}
