package main

import (
	"bufio"
	"fmt"
	"io"
	"net"
)

func writeConn(conn net.Conn, messageChan chan string) {
	for {
		msg, ok := <-messageChan
		if !ok {
			return
		}
		fmt.Println(msg)
		_, _ = io.WriteString(conn, msg)
	}
}

func handleConn(conn net.Conn) {
	defer conn.Close()
	messageChan := make(chan string)

	go writeConn(conn, messageChan)

	reader := bufio.NewReader(conn)
	for {
		msg, err := reader.ReadString('\n')
		if err != nil {
			close(messageChan)
			return
		}
		messageChan <- msg
	}
}

func main() {
	listener, err := net.Listen("tcp", "127.0.0.1:8080")
	if err != nil {
		panic(err)
	}

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println(err)
			return
		}

		go handleConn(conn)
	}
}
