package main

import (
	"fmt"
	"net"
	"sync"
)

type Server struct {
	Ip   string
	Port int

	// 已连接到服务器的用户列表
	UserMap map[string]*User
	MapLock sync.RWMutex

	BroadcastChan chan string
}

func InitServer(ip string, port int) *Server {
	server := &Server{
		Ip:            ip,
		Port:          port,
		UserMap:       make(map[string]*User),
		BroadcastChan: make(chan string),
	}
	return server
}

func (server *Server) Start() {
	// socket listen
	listener, err := net.Listen("tcp", fmt.Sprintf("%s:%d", server.Ip, server.Port))
	if err != nil {
		fmt.Printf("server-1 socket listen err:%s \n", err)
		return
	}
	// socket listener close at last
	defer func(listener net.Listener) {
		err := listener.Close()
		if err != nil {
			fmt.Printf("server-1 socket listener close err:%s \n", err)
		}
	}(listener)

	go server.StartBroadcast()

	for {
		// accept connection
		connection, err := listener.Accept()
		if err != nil {
			fmt.Printf("server-1 socket accept err:%s \n", err)
			continue
		}
		// do handler
		go server.Handler(connection)
	}
}

func (server *Server) Handler(conn net.Conn) {
	// 创建当前用户对象
	user := InitUser(conn)

	// 将其维护到用户表中
	server.MapLock.Lock()
	server.UserMap[user.Name] = user
	server.MapLock.Unlock()

	// 广播发送连接成功的消息
	msg := fmt.Sprintf("连接到 server-1:%s 成功", fmt.Sprintf("%s:%d", server.Ip, server.Port))
	server.SendMessage(user, msg)

	// 阻塞当前函数
	//		主要作用为：防止 user 等临时变量被销毁，UserMap中的引用失效
	select {}
}

func (server *Server) SendMessage(user *User, msg string) {
	broadcastMsg := "[" + user.Addr + "]" + user.Name + ":" + msg
	server.BroadcastChan <- broadcastMsg
}

func (server *Server) StartBroadcast() {
	for {
		// 获取待发送的广播信息
		msg := <-server.BroadcastChan

		// 将广播信息发给所有已连接的用户
		server.MapLock.Lock()
		for _, user := range server.UserMap {
			user.ReceiveChan <- msg
		}
		server.MapLock.Unlock()
	}

}
