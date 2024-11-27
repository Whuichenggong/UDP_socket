package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
)

func main() {
	//1. 创建UDP连接到服务器的地址和端口号
	c, err := net.DialUDP("udp", nil, &net.UDPAddr{
		IP:   net.IPv4(127, 0, 0, 1),
		Port: 8282,
	})
	if err != nil {
		fmt.Println("dial err: %v\n", err)
		return
	}
	defer c.Close() // 将 defer 放在 if 语句外面

	// 2.从标准输入读取用户输入的数据
	input := bufio.NewReader(os.Stdin)
	for {
		// 读取用户输入知道遇见换行符
		s, err := input.ReadString('\n')
		if err != nil {
			fmt.Printf("read from stdin failed, err: %v\n", err)
			return
		}

		//3. 将用户输入的数据转换为字节数组并通过UDP连接发送给服务器
		_, err = c.Write([]byte(s))
		if err != nil {
			fmt.Printf("send to server failed, err: %v\n", err)
			return
		}

		// 4.接收来自服务器的数据
		var buf [1024]byte
		n, addr, err := c.ReadFromUDP(buf[:])
		if err != nil {
			fmt.Printf("recv from udp failed, err: %v\n", err)
			return
		}

		// 打印来自服务器的数据
		fmt.Printf("服务器 %v, 响应数据: %v\n", addr, string(buf[:n]))
	}
}
