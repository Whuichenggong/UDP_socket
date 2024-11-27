package main

import (
	"fmt"
	"net"
)

// udp server
func main() {
	// 创建一个UDP监听器，监听本地IP地址的端口
	listen, err := net.ListenUDP("udp", &net.UDPAddr{
		IP:   net.IPv4(127, 0, 0, 1),
		Port: 8282,
	})
	if err != nil {
		fmt.Printf("listen failed,err:%v\n", err)
		return
	}
	defer listen.Close()

	for {
		var buf [1024]byte
		// 从UDP连接中读取数据到buf中，n为读取到的字节数，addr为数据发送者的地址
		n, addr, err := listen.ReadFromUDP(buf[:])
		if err != nil {
			fmt.Printf("read from udp failed,err:%v\n", err)
			return
		}

		// 打印接收到的数据
		fmt.Println("接收到的数据：", string(buf[:n]))

		// 将接收到的数据原样发送回给数据发送者
		_, err = listen.WriteToUDP(buf[:n], addr)
		if err != nil {
			fmt.Printf("write to %v failed,err:%v\n", addr, err)
			return
		}
	}
}
