package main

import (
	"flag"
	"fmt"
	"io"
	"net"
	"os"
)

type Client struct {
	ServerIp   string
	ServerPort int
	Name       string
	conn       net.Conn
	flag       int
}

func NewClient(serverIp string, serverPort int) *Client {
	client := &Client{
		ServerIp:   serverIp,
		ServerPort: serverPort,
		flag:       999,
	}
	conn, err := net.Dial("tcp", fmt.Sprintf("%s:%d", serverIp, serverPort))
	if err != nil {
		fmt.Println("net.Dial error:", err)
		return nil
	}
	client.conn = conn
	return client
}
func (client *Client) DealResponse() {
	io.Copy(os.Stdout, client.conn)
}
func (client *Client) menu() bool {
	fmt.Println("1. 公聊模式")
	fmt.Println("2. 私聊模式")
	fmt.Println("3. 更新用户名")
	fmt.Println("0. 退出")
	var flag int
	fmt.Scanln(&flag)
	if flag >= 0 && flag <= 3 {
		client.flag = flag
		return true
	}

	fmt.Println("请输入合法范围内的数字...")
	return false
}
func (client *Client) UpdateName() bool {
	fmt.Println(">>>>请输入用户名：")
	fmt.Scanln(&client.Name)
	sendMsg := "rename|" + client.Name + "\n"
	_, err := client.conn.Write([]byte(sendMsg))
	if err != nil {
		fmt.Println("conn.Write error:", err)
		return false
	}
	return true
}
func (client Client) PublicChat() {
	fmt.Println(">>>>请输入聊天内容，exit退出")
	var chatMsg string
	fmt.Scanln(&chatMsg)
	for chatMsg != "exit" {
		if len(chatMsg) != 0 {
			sendMsg := chatMsg + "\n"
			_, err := client.conn.Write([]byte(sendMsg))
			if err != nil {
				fmt.Println("conn.Write error:", err)
				return
			}
		}
		chatMsg = ""
		fmt.Scanln(&chatMsg)
	}

}
func (client Client) PrivateChat() {
	client.SelectUsers()
	fmt.Println(">>>>请输入聊天对象[用户名], exit退出")
	var remoteName string
	fmt.Scanln(&remoteName)
	for remoteName != "exit" {
		fmt.Println(">>>>请输入消息内容，exit退出")
		var chatMsg string
		fmt.Scanln(&chatMsg)
		for chatMsg != "exit" {
			if len(chatMsg) != 0 {
				sendMsg := "to|" + remoteName + "|" + chatMsg + "\n"
				_, err := client.conn.Write([]byte(sendMsg))
				if err != nil {
					fmt.Println("conn.Write error:", err)
					return
				}
			}
			chatMsg = ""
			fmt.Scanln(&chatMsg)
		}
		remoteName = ""
		fmt.Scanln(&remoteName)
	}

}
func (client Client) SelectUsers() {
	sendMsg := "who\n"
	_, err := client.conn.Write([]byte(sendMsg))
	if err != nil {
		fmt.Println("conn.Write error:", err)
		return
	}
}
func (client *Client) Run() {
	for client.flag != 0 {
		for client.menu() == false {

		}
		switch client.flag {
		case 1:
			client.PublicChat()
			break
		case 2:
			client.PrivateChat()
			break
		case 3:
			client.SelectUsers()
			break
		}
	}
}

var serverIp string
var serverPort int

func init() {
	flag.StringVar(&serverIp, "serverIp", "127.0.0.1", "server ip")
	flag.IntVar(&serverPort, "serverPort", 8888, "server port")
}
func main() {
	flag.Parse()
	client := NewClient(serverIp, serverPort)
	if client == nil {
		fmt.Println(">>>>> 连接服务器失败...")
		return
	}
	go client.DealResponse()
	defer client.conn.Close()
	fmt.Println(">>>>> 连接服务器成功...")

	client.Run()
}
