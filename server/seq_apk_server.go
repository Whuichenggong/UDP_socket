package main

import (
	"fmt"
	"net"
	"strconv"
	"strings"
)

type Message struct {
	Seq int
	Msg string
}

func main() {
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
		n, addr, err := listen.ReadFromUDP(buf[:])
		if err != nil {
			fmt.Printf("read from udp failed,err:%v\n", err)
			return
		}

		// 处理接收到的数据，提取序列号和消息文本
		message := decodeMessage(buf[:n])
		fmt.Printf("Received seq=%d from %v: %s\n", message.Seq, addr, message.Msg)

		// 发送ACK回复给客户端，ACK=Seq+1
		ack := Message{Seq: message.Seq + 1, Msg: "ACK"}
		_, err = listen.WriteToUDP(encodeMessage(ack), addr)
		if err != nil {
			fmt.Printf("write to %v failed,err:%v\n", addr, err)
			return
		}
	}
}

func encodeMessage(msg Message) []byte {
	// 将序列号和消息文本编码成字节数据
	return []byte(fmt.Sprintf("%d;%s", msg.Seq, msg.Msg))
}

func decodeMessage(data []byte) Message {
	// 解码收到的数据，提取序列号和消息文本
	parts := strings.Split(string(data), ";")
	seq, _ := strconv.Atoi(parts[0])
	msg := parts[1]
	return Message{Seq: seq, Msg: msg}
}
