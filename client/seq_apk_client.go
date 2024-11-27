package main

import (
	"fmt"
	"net"
	"strconv"
	"strings"
	"time"
)

type Message struct {
	Seq int
	Msg string
}

func main() {
	c, err := net.DialUDP("udp", nil, &net.UDPAddr{
		IP:   net.IPv4(127, 0, 0, 1),
		Port: 8282,
	})
	if err != nil {
		fmt.Printf("dail,err:%v\n", err)
		return
	}
	defer c.Close()

	//示例数据
	input := []string{"Message 1", "Message 2", "Message 3", "Message 4", "Message 5"}
	seq := 0

	for _, msg := range input {
		seq++
		message := Message{Seq: seq, Msg: msg}
		fmt.Printf("Sending seq=%d: %s\n", message.Seq, message.Msg)

		// 发送带有序列号的数据包
		_, err = c.Write(EncodeMessage(message))
		if err != nil {
			fmt.Printf("send to server failed,err:%v\n", err)
			return
		}

		// 等待ACK，设置超时时间
		buf := make([]byte, 1024)
		c.SetReadDeadline(time.Now().Add(5 * time.Second))
		n, _, err := c.ReadFromUDP(buf)
		if err != nil {
			fmt.Println("ACK not received. Timeout or Error.")
			return
		} else {
			ack := DecodeMessage(buf[:n])
			if ack.Seq == seq+1 {
				fmt.Printf("ACK = %d\n", ack.Seq)
			} else {
				fmt.Println("Invalid ACK received. Retry.")
				return
			}
		}
	}
}

func EncodeMessage(msg Message) []byte {
	// 将序列号和消息文本编码成字节数据
	return []byte(fmt.Sprintf("%d;%s", msg.Seq, msg.Msg))
}

func DecodeMessage(data []byte) Message {
	// 解码收到的数据，提取序列号和消息文本
	parts := strings.Split(string(data), ";")
	seq, _ := strconv.Atoi(parts[0])
	msg := parts[1]
	return Message{Seq: seq, Msg: msg}
}
