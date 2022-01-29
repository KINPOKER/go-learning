package main

import (
	"flag"
	"fmt"
	"io"
	"net"
	"os"
	"time"
)

type Client struct {
	ServerIp   string
	ServerPort int
	Name       string
	Conn       net.Conn
	Mode       int
}

func InitClient(clientName string, serverIp string, serverPort int) *Client {
	client := &Client{
		ServerIp:   serverIp,
		ServerPort: serverPort,
		Name:       clientName,
		Mode:       -1,
	}

	dialConn, err := net.Dial("tcp", fmt.Sprintf("%s:%d", serverIp, serverPort))
	if err != nil {
		fmt.Printf("client dial err:%s \n", err)
		return nil
	}

	client.Conn = dialConn
	return client
}

func (client *Client) ShowMenu() bool {
	var clientMode int

	fmt.Println("1:公聊模式")
	fmt.Println("2:私聊模式")
	fmt.Println("3:更换用户名")
	fmt.Println("0:退出")

	_, err := fmt.Scanln(&clientMode)
	if err != nil {
		fmt.Println("请输入合法的数字，以选择客户端模式")
		return false
	}

	if clientMode > -1 && clientMode < 4 {
		client.Mode = clientMode
		return true
	} else {
		fmt.Println("请输入合法的数字，以选择客户端模式")
		return false
	}
}

func (client *Client) Start() {
	for client.Mode != 0 {
		for client.ShowMenu() != true {
		}

		switch client.Mode {
		case 1:
			fmt.Println("客户端进入公聊模式")
			client.StartPublicChat()
			break
		case 2:
			fmt.Println("客户端进入私聊模式")
			client.StartPrivateChat()
			break
		case 3:
			fmt.Println("可以开始更换用户名")
			client.UpdateClientName()
			break
		}
	}
}

func (client *Client) StartPublicChat() {
	var publicChatMsg string
	fmt.Printf("请输入待发送的消息，输入 %s 退出公聊模式。\n", exitPublicChatKey)

	_, err := fmt.Scanln(&publicChatMsg)
	if err != nil {
		fmt.Println("请输入合法的字符串作为待发送的消息")
		return
	}

	for publicChatMsg != exitPublicChatKey {
		_, err = client.Conn.Write([]byte(publicChatMsg + "\n"))
		if err != nil {
			fmt.Printf("client conn write error:%s \n", err)
			return
		}

		time.Sleep(time.Millisecond * time.Duration(userInputTimeDuration))

		fmt.Printf("请输入待发送的消息，输入 %s 退出公聊模式。\n", exitPublicChatKey)
		_, err = fmt.Scanln(&publicChatMsg)
		if err != nil {
			fmt.Println("请输入合法的字符串作为待发送的消息")
			return
		}
	}
}

func (client *Client) StartPrivateChat() {
	var privateChatMsg string
	var sendToUserName string

	fmt.Printf("所有在线用户如下，请输入待私聊的用户名，输入 %s 退出私聊模式。\n", exitPrivateChatUserChooseKey)
	client.GetUsersOline()
	_, err := fmt.Scanln(&sendToUserName)
	if err != nil {
		fmt.Println("请输入合法的字符串作为待私聊的用户名")
		return
	}
	for sendToUserName != exitPrivateChatUserChooseKey {
		fmt.Printf("请输入待发送的消息，输入 %s 退出与 %s 的私聊。\n", exitPrivateChatKey, sendToUserName)
		_, err = fmt.Scanln(&privateChatMsg)
		if err != nil {
			fmt.Println("请输入合法的字符串作为待发送的消息")
			return
		}

		for privateChatMsg != exitPrivateChatKey {
			msg := privateChatKey + sendToUserName + "|" + privateChatMsg
			_, err = client.Conn.Write([]byte(msg + "\n"))
			if err != nil {
				fmt.Printf("client conn write error:%s \n", err)
				return
			}

			time.Sleep(time.Millisecond * time.Duration(userInputTimeDuration))

			fmt.Printf("请输入待发送的消息，输入 %s 退出与 %s 的私聊。\n", exitPrivateChatKey, sendToUserName)
			_, err = fmt.Scanln(&privateChatMsg)
			if err != nil {
				fmt.Println("请输入合法的字符串作为待发送的消息")
				return
			}
		}

		time.Sleep(time.Millisecond * time.Duration(userInputTimeDuration))

		fmt.Printf("所有在线用户如下，请输入待私聊的用户名，输入 %s 退出私聊模式。\n", exitPrivateChatUserChooseKey)
		client.GetUsersOline()
		_, err = fmt.Scanln(&sendToUserName)
		if err != nil {
			fmt.Println("请输入合法的字符串作为待私聊的用户名")
			return
		}
	}

}

func (client *Client) GetUsersOline() {
	msg := onlineUserQueryKey + "\n"
	_, err := client.Conn.Write([]byte(msg))
	if err != nil {
		fmt.Printf("client conn write error:%s \n", err)
	}
	return
}

func (client *Client) UpdateClientName() {
	var newClientName string
	fmt.Println("请输入新用户名：")

	_, err := fmt.Scanln(&newClientName)
	if err != nil {
		fmt.Println("请输入合法的字符串作为新用户名")
		return
	}

	msg := renameFuncKey + newClientName + "\n"
	_, err = client.Conn.Write([]byte(msg))
	if err != nil {
		fmt.Printf("client conn write error:%s \n", err)
	}
	return
}

func (client *Client) MessageSync() {
	_, err := io.Copy(os.Stdout, client.Conn)
	if err != nil {
		fmt.Printf("client message from io copy error:%s \n", err)
		return
	}
}

var clientName string
var serverIp string
var serverPort int

func init() {
	flag.StringVar(&clientName, "clientName", "client1", "设置客户端名称")
	flag.StringVar(&serverIp, "serverIp", "127.0.0.1", "设置服务端 ip")
	flag.IntVar(&serverPort, "serverPort", 8888, "设置服务端端口号")
}

func main() {
	// 命令行解析
	flag.Parse()

	client := InitClient(clientName, serverIp, serverPort)
	if client == nil {
		fmt.Printf("客户端连接到 server 失败 \n")
	} else {
		fmt.Printf("客户端 %s 连接到 server 成功 \n", client.Name)
	}

	go client.MessageSync()

	client.Start()
}
