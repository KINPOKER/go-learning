package main

import (
	"fmt"
	"net"
	"strings"
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
	for msg := range user.ReceiveChan {
		_, err := user.Conn.Write([]byte(msg))
		if err != nil {
			fmt.Printf("用户:%s 客户端打印消息失败,消息：%s,error:%s \n", user.Name, msg, err)
			return
		}
	}
}

func (user *User) login(server *Server) {
	// 将其维护到用户表中
	server.MapLock.Lock()
	server.UserMap[user.Name] = user
	server.MapLock.Unlock()

	// 广播发送连接成功的消息
	msg := fmt.Sprintf("连接到 server:%s 成功", fmt.Sprintf("%s:%d", server.Ip, server.Port))
	server.SendMessage(user, msg)
}

func (user *User) logout(server *Server) {
	// 同步维护用户表
	server.MapLock.Lock()
	delete(server.UserMap, user.Name)
	server.MapLock.Unlock()

	server.SendMessage(user, "断开连接")
}

func (user *User) handleMessage(server *Server, msg string) {
	onlineUserQueryKey := "who is online?"
	renameFuncKey := "rename|"

	if msg == onlineUserQueryKey {
		server.MapLock.Lock()
		for _, curUser := range server.UserMap {
			user.printMessage("[" + curUser.Addr + "]" + curUser.Name + "在线\n")
		}
		server.MapLock.Unlock()
	} else if len(msg) > len(renameFuncKey) && strings.HasPrefix(msg, renameFuncKey) {
		newName := msg[len(renameFuncKey):]

		// 同步维护用户表
		server.MapLock.Lock()
		_, ok := server.UserMap[newName]
		if ok {
			user.printMessage(fmt.Sprintf("用户名:%s已被使用,请更换后重试\n", newName))
		} else {
			delete(server.UserMap, user.Name)
			user.Name = newName
			server.UserMap[newName] = user
			user.printMessage(fmt.Sprintf("用户名成功更换为:%s \n", newName))
		}
		server.MapLock.Unlock()
	} else {
		server.SendMessage(user, msg)
	}
}

func (user *User) printMessage(msg string) {
	_, err := user.Conn.Write([]byte(msg))
	if err != nil {
		fmt.Printf("用户:%s 客户端打印消息失败,消息：%s,error:%s \n", user.Name, msg, err)
		return
	}
}
