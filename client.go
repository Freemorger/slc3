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
	const HOST_NAME = "localhost"
	const HOST_PORT = 8081
	const CONN_TYPE = "tcp"
	var domain = (HOST_NAME + ":" + strconv.Itoa(HOST_PORT))

	conn, err := net.Dial(CONN_TYPE, domain)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	fmt.Println("Connected to server: ", domain)
	defer conn.Close()
	go ConReader(conn)
	fmt.Print("You > ")
	for {
		consoleScanner := bufio.NewScanner(os.Stdin)
		for consoleScanner.Scan() {
			fmt.Print("You > ")
			text := consoleScanner.Text()
			if strings.ToLower(text) == "exit" {
				os.Exit(0)
			}
			_, err := conn.Write([]byte(text + "\n"))
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

func ConReader(conn net.Conn) {
	for {
		response, err := bufio.NewReader(conn).ReadString('\n')
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		fmt.Println()
		fmt.Print(response)
	}

}
