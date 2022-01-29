package main

import (
	"fmt"
	"io"
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
	user.login(server)

	// 接受用户发送的消息
	go func() {
		buf := make([]byte, 4096)
		for {
			length, err := conn.Read(buf)
			if length == 0 {
				user.logout(server)
				return
			}

			if err != nil && err != io.EOF {
				fmt.Printf("Conn read error:%s \n", err)
				fmt.Printf("server读取user:%s 发送的信息失败 \n", user.Name)
				return
			}

			// 提取用户发送的消息（去掉换行符）
			msg := string(buf[:length-1])
			user.handleMessage(server, msg)
		}
	}()

	// 阻塞当前函数
	//		主要作用为：防止 user 等临时变量被销毁，UserMap中的引用失效
	select {}
}

func (server *Server) SendMessage(user *User, msg string) {
	broadcastMsg := "[" + user.Addr + "]" + user.Name + ":" + msg + " \n"
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
