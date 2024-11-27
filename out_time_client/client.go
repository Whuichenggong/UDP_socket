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
		fmt.Printf("dial,err:%v\n", err)
		return
	}
	defer c.Close()

	// 示例数据
	input := []string{"Message 1", "Message 2", "Message 3", "Message 4", "Message 5"}

	for seq, msg := range input {
		for {
			message := Message{Seq: seq + 1, Msg: msg}
			fmt.Printf("Sending seq=%d: %s\n", message.Seq, message.Msg)

			// 发送带有序列号的数据包
			_, err := c.Write(encodeMessage(message))
			if err != nil {
				fmt.Printf("send to server failed,err:%v\n", err)
				return
			}

			// 开始等待ACK，设置超时时间
			buf := make([]byte, 1024)
			c.SetReadDeadline(time.Now().Add(5 * time.Second))

			// 循环等待ACK，直到收到正确的ACK或超时
			n, _, err := c.ReadFromUDP(buf)
			if err != nil {
				// 超时或发生错误，需要重传
				fmt.Println("ACK not received. Timeout or Error. Retrying...")
				continue
			} else {
				//解码从服务器传来的ack
				ack := decodeMessage(buf[:n])
				if ack.Seq == seq+2 {
					fmt.Printf("ACK = %d\n", ack.Seq)
					// 收到正确的ACK，跳出内部循环，继续发送下一个消息
					break
				} else {
					// 收到错误的ACK，继续等待，内部循环会重发相同的消息
					fmt.Println("Invalid ACK received. Waiting for correct ACK...")
					continue
				}
			}
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
