package main

import (
	"fmt"
	"net"
)

type Server struct {
	Ip   string
	Port int
}

func InitServer(ip string, port int) *Server {
	server := &Server{
		Ip:   ip,
		Port: port,
	}
	return server
}

func (server *Server) Start() {
	// socket listen
	listener, err := net.Listen("tcp", fmt.Sprintf("%s:%d", server.Ip, server.Port))
	if err != nil {
		fmt.Printf("server socket listen err:%s\n", err)
		return
	}
	// socket listener close at last
	defer listener.Close()

	for {
		// accept connection
		connection, err := listener.Accept()
		if err != nil {
			fmt.Printf("server socket accept err:%s\n", err)
			continue
		}
		// do handler
		go server.Handler(connection)
	}
}

func (server *Server) Handler(conn net.Conn) {
	fmt.Printf("连接到 server:%s 成功", fmt.Sprintf("%s:%d", server.Ip, server.Port))
}
