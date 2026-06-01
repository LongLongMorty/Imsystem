package main

import (
	"fmt"
	"net"
)

type Server struct {
	Ip   string
	Port int
}

// 创建Server接口
func NewServer(ip string, port int) *Server {
	return &Server{
		Ip:   ip,
		Port: port,
	}
}

// 启动服务器的接口
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
	fmt.Println("链接建立成功")
}
