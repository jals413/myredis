package main

import (
	"fmt"
	"net"
	"os"
	"strconv"
	"strings"
)

func main() {
	l, err := net.Listen("tcp", "0.0.0.0:6379")
	if err != nil {
		fmt.Println("Failed to bind to port 6379")
		os.Exit(1)
	}
	defer l.Close()
	for {
		conn, err := l.Accept()
		if err != nil {
			fmt.Println("Error accepting connection: ", err.Error())
			continue
		}
		go handleConnection(conn)
	}
}

func handleConnection(conn net.Conn) {
	defer conn.Close()

	buf := make([]byte, 1024)
	for {
		cmd, err := conn.Read(buf)
		if err != nil {
			fmt.Println("Error reading from client: ", err.Error())
			return
		}
		handleCommands(string(buf[:cmd]), conn)
	}
}

func handleCommands(s string, conn net.Conn) {

	lines := strings.Split(s, "\r\n") 
	lines = lines[:len(lines)-1] // removing last empyt character.

	if lines[0][0] != '*' {
		conn.Write([]byte("-invalid command\r\n"))
	}

	num, err := strconv.Atoi(lines[0][1:])
	if err != nil || 2*num+1 != len(lines) || num > 2 {
		conn.Write([]byte("-invalid command\r\n"))
	}

	switch strings.ToUpper(lines[2]) {
	case "PING":
		conn.Write([]byte("+PONG\r\n"))
	case "ECHO":
		if len(lines) != 5 {
			conn.Write([]byte("-invalid command\r\n"))
		}
		conn.Write([]byte("+" + lines[4] + "\r\n"))
	default:
		conn.Write([]byte("-unknown command\r\n"))
	}

} 
