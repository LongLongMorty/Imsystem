package main

import "net"

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
		msg := <-u.C
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

	} else {
		u.server.Broadcast(u, msg)
	}
}

// SendMsg 给当前用户发送消息
func (u *User) SendMsg(msg string) {
	u.conn.Write([]byte(msg))

}
