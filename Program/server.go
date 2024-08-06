package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strconv"
	"strings"
)

func main() {
	const NAME = "localhost"
	const PORT = 8081
	const CONN_TYPE = "tcp"
	var domain = (NAME + ":" + strconv.Itoa(PORT))

	listener, err := net.Listen(CONN_TYPE, domain)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer listener.Close()
	fmt.Println(CONN_TYPE, " Listener started on ", domain)

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		fmt.Println("New connection estabilished: ", conn.RemoteAddr().String())
		go handleRequest(conn)
		fmt.Print("Server > ")
		for {
			consoleScanner := bufio.NewScanner(os.Stdin)
			for consoleScanner.Scan() {
				fmt.Print("Server > ")
				text := consoleScanner.Text()
				if strings.ToLower(text) == "exit" {
					conn.Write([]byte("Server is closing."))
					conn.Close()
					os.Exit(0)
				}
				_, err := conn.Write([]byte("SERVER: " + text + "\n"))
				if err != nil {
					fmt.Println(err)
					os.Exit(1)
				}
			}

			if err := consoleScanner.Err(); err != nil {
				fmt.Println("Error reading from terminal: ", err.Error())
			}
		}
	}
}

func handleRequest(conn net.Conn) {
	defer conn.Close()
	scanner := bufio.NewScanner(conn)
	for scanner.Scan() {
		clientMessage := scanner.Text()
		fmt.Println(conn.RemoteAddr(), ": ", clientMessage)
		_, err := conn.Write([]byte(conn.RemoteAddr().String() + ": " + clientMessage))
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	}

	if err := scanner.Err(); err != nil {
		fmt.Println("Error reading: ", err.Error())
	}
}
