package main

import "fmt"

// 协程：往通道发数据
func send(ch chan int) {
	fmt.Println("协程正在放数据：42")
	ch <- 42 // 放数据 → 阻塞等待有人取
	fmt.Println("数据被取走了！")
}

func main() {
	// 1. 创建通道
	ch := make(chan int)

	// 2. 启动协程
	go send(ch)

	// 3. 主协程从通道取数据
	fmt.Println("主协程等待取数据...")
	num := <-ch // 这里会阻塞，直到有人放数据

	fmt.Println("取到数据：", num)
}
