package main

import (
	"bufio"
	"fmt"
	"net"
)

var clients = make(map[net.Conn]bool)

func main() {
	listener, _ := net.Listen("tcp", ":8080")
	fmt.Println("Сервер на :8080")
	
	for {
		conn, _ := listener.Accept()
		clients[conn] = true
		go handle(conn)
	}
}

func handle(conn net.Conn) {
	defer conn.Close()
	
	reader := bufio.NewReader(conn)
	for {
		msg, _ := reader.ReadString('\n')
		if msg == "" {
			return
		}
		
		// Отправляем всем
		for c := range clients {
			if c != conn {
				c.Write([]byte(msg))
			}
		}
		fmt.Print(msg)
	}
}
