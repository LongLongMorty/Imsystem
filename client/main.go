package main

import (
	"bufio"
	"fmt"
	"io"
	"net"
	"os"
)

func main() {
	conn, err := net.Dial("tcp", "127.0.0.1:8888")
	if err != nil {
		fmt.Println("连接服务器失败:", err)
		return
	}
	defer conn.Close()

	fmt.Println("连接服务器成功！你可以开始输入消息了 (按 Ctrl+C 退出)")

	// 启动一个 goroutine，专门负责接收服务端的广播消息并打印到终端
	go func() {
		// 将连接中传来的数据直接拷贝到标准输出
		io.Copy(os.Stdout, conn)
	}()

	// 阻塞读取用户在当前窗口的键盘输入，并发送给服务端
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		text := scanner.Text()
		_, err := conn.Write([]byte(text + "\n"))
		if err != nil {
			fmt.Println("发送失败:", err)
			break
		}
	}
}
