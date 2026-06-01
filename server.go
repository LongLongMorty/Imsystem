package main

import (
	"fmt"
	"net"
	"sync"
)

type Server struct {
	Ip        string
	Port      int
	OnlineMap map[string]*User
	mapLock   sync.RWMutex
	Message   chan string
}

// NewServer 创建Server接口
func NewServer(ip string, port int) *Server {
	return &Server{
		Ip:        ip,
		Port:      port,
		OnlineMap: make(map[string]*User),
		Message:   make(chan string),
	}
}

// Start 启动服务器的接口
func (s *Server) Start() {
	// socekt.Listen()
	listener, err := net.Listen("tcp", fmt.Sprintf("%s:%d", s.Ip, s.Port))
	if err != nil {
		fmt.Println("net.Listen error:", err)
		return
	}
	fmt.Println("服务器启动成功，等待链接...")
	// close listen socket
	defer listener.Close()

	// 启动监听Message的goroutine
	go s.ListenMassage()

	for {
		// accept
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("listener.Accept error:", err)
			continue
		}
		// do handler
		go s.Handler(conn)
	}

}
func (s *Server) Handler(conn net.Conn) {
	//用户上线，将用户加入OnlineMap中
	user := NewUser(conn)
	s.mapLock.Lock()
	s.OnlineMap[user.Name] = user
	s.mapLock.Unlock()

	//广播当前用户上线消息
	s.Broadcast(user, "已上线")
	// 接受客户端发送的消息
	go func() {
		buf := make([]byte, 4096)
		for {
			n, err := conn.Read(buf)
			if err != nil {
				s.Broadcast(user, "下线了")
				s.mapLock.Lock()
				delete(s.OnlineMap, user.Name)
				s.mapLock.Unlock()
				break
			}
			msg := string(buf[:n-1]) // 去掉换行符
			s.Broadcast(user, msg)
		}
	}()
	// 阻塞客户端
	select {}

}
func (s *Server) Broadcast(user *User, msg string) {
	sendMsg := "[" + user.Addr + "]" + user.Name + ":" + msg
	s.Message <- sendMsg
}

// ListenMassage 监听Message广播消息的goroutine，一旦有消息就发送给全部在线User
func (s *Server) ListenMassage() {
	for {
		msg := <-s.Message
		s.mapLock.RLock()
		for _, u := range s.OnlineMap {
			u.C <- msg
		}
		s.mapLock.RUnlock()
	}
}
