package main

import (
	"net"
	"strings"
)

type User struct {
	Name   string
	Addr   string
	C      chan string
	conn   net.Conn
	server *Server
}

// NewUser 创建用户的API
func NewUser(conn net.Conn, server *Server) *User {
	userAddr := conn.RemoteAddr().String()
	user := &User{
		Name:   userAddr,
		Addr:   userAddr,
		C:      make(chan string),
		conn:   conn,
		server: server,
	}
	go user.ListenMessage()
	return user
}

// ListenMessage 监听当前User channel的方法，一旦有消息，就直接发送给对端客户端
func (u *User) ListenMessage() {
	for {
		msg, ok := <-u.C
		if !ok {
			break
		}
		u.conn.Write([]byte(msg + "\n"))
	}
}

// Online 用户上线
func (u *User) Online() {
	//用户上线，将用户加入OnlineMap中
	u.server.mapLock.Lock()
	u.server.OnlineMap[u.Name] = u
	u.server.mapLock.Unlock()
	//广播当前用户上线消息
	u.server.Broadcast(u, "已上线")
}

// Offline 用户下线
func (u *User) Offline() {
	//用户下线，将用户加入OnlineMap中
	u.server.mapLock.Lock()
	delete(u.server.OnlineMap, u.Name)
	u.server.mapLock.Unlock()
	//广播当前用户下线消息
	u.server.Broadcast(u, "已下线")
}

// DoMessage 用户端发送消息时，执行DoMessage方法
func (u *User) DoMessage(msg string) {
	if msg == "who" {
		// 查询当前在线用户都有哪些
		u.server.mapLock.RLock()
		for _, user := range u.server.OnlineMap {
			onlineMsg := "[" + user.Addr + "]" + user.Name + ":" + "在线...\n"
			u.SendMsg(onlineMsg)
		}
		u.server.mapLock.RUnlock()

	} else if len(msg) > 7 && msg[:7] == "rename|" {
		// 消息格式：rename|张三
		newName := msg[7:]
		// 判断name是否存在
		_, ok := u.server.OnlineMap[newName]
		if ok {
			u.SendMsg("当前用户名被占用\n")
		} else {
			u.server.mapLock.Lock()
			delete(u.server.OnlineMap, u.Name) // 删除原来的key
			u.Name = newName
			u.server.OnlineMap[u.Name] = u // 以新的key加入map
			u.server.mapLock.Unlock()
			u.SendMsg("您已经更新用户名:" + u.Name + "\n")
		}

	} else if len(msg) > 3 && msg[:3] == "to|" {
		// 消息格式：to|张三|消息内容
		// 1. 获取对方用户名和消息内容
		parts := strings.Split(msg[3:], "|")
		if len(parts) < 2 {
			u.SendMsg("消息格式不正确，请使用\"to|张三|消息内容\"格式\n")
			return
		}
		remoteName := parts[0]
		content := parts[1]

		if remoteName == "" || content == "" {
			u.SendMsg("消息格式不正确，请使用\"to|张三|消息内容\"格式\n")
			return
		}
		// 2. 查找对方用户
		u.server.mapLock.RLock()
		remoteUser, ok := u.server.OnlineMap[remoteName]
		u.server.mapLock.RUnlock()
		if !ok {
			u.SendMsg("该用户名不存在\n")
			return
		}
		if remoteUser == u {
			u.SendMsg("不能给自己发送消息\n")
			return
		}
		// 3. 发送消息
		remoteUser.SendMsg(u.Name + "对您说:" + content + "\n")

	} else {
		u.server.Broadcast(u, msg)
	}
}

// SendMsg 给当前用户发送消息
func (u *User) SendMsg(msg string) {
	u.conn.Write([]byte(msg))

}
