package main

import (
	"fmt"
	"net"
)

type Client struct {
	ServerIp   string
	ServerPort int
	Name       string
	Conn       net.Conn
}

func InitClient(clientName string, serverIp string, serverPort int) *Client {
	client := &Client{
		ServerIp:   serverIp,
		ServerPort: serverPort,
		Name:       clientName,
	}

	dialConn, err := net.Dial("tcp", fmt.Sprintf("%s:%d", serverIp, serverPort))
	if err != nil {
		fmt.Printf("client dial err:%s \n", err)
		return nil
	}

	client.Conn = dialConn
	return client
}

func main() {
	client := InitClient("client1", "127.0.0.1", 8888)
	if client == nil {
		fmt.Printf("客户端连接到 server 失败 \n")
	} else {
		fmt.Printf("客户端%s连接到 server 成功 \n", client.Name)
	}

	select {}
}
