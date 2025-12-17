package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
)

func main() {
	conn, _ := net.Dial("tcp", "localhost:8080")
	
	// Читаем с сервера
	go func() {
		reader := bufio.NewReader(conn)
		for {
			msg, _ := reader.ReadString('\n')
			fmt.Print(msg)
		}
	}()
	
	// Читаем с клавиатуры и отправляем
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		text := scanner.Text()
		conn.Write([]byte(text + "\n"))
	}
}
