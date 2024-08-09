package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strconv"
	"strings"
	"time"
)

type ClientSelf struct {
	connection net.Conn
	name       string
}

func main() {
	// Server info for connection
	const HOST_NAME = "localhost"
	const HOST_PORT = 8081
	const CONN_TYPE = "tcp"
	var clientName string
	var client ClientSelf

	fmt.Print("Enter your name: ")
	fmt.Scanln(&clientName)
	client = connectToServer(HOST_NAME, HOST_PORT, CONN_TYPE, clientName)
	defer disconnect(client)

	go ConReader(client)
	ConWriter(client)

}

func connectToServer(hostname string, port int, conntype string,
	clientname string) (client ClientSelf) {
	var domain = (hostname + ":" + strconv.Itoa(port))
	// Connect to server
	connection, err := net.Dial(conntype, domain)
	if err != nil {
		fmt.Println("Error connecting to server: " + err.Error())
		os.Exit(1)
	}
	connection.Write([]byte(clientname))
	fmt.Println("Connected to server: ", domain)
	clientself := ClientSelf{connection: connection, name: clientname}
	return clientself
}

func ConWriter(client ClientSelf) {
	consoleScanner := bufio.NewScanner(os.Stdin)
	fmt.Print("You > ")
	// Loop for reading user input
	for {
		for consoleScanner.Scan() {
			fmt.Print("You > ")
			text := consoleScanner.Text()
			if strings.ToLower(text) == "exit" {
				disconnect(client)
			}
			_, err := client.connection.Write([]byte(text + "\n"))
			if err != nil {
				fmt.Println("Sending message error: ", err.Error())
				continue
			}
		}

		if err := consoleScanner.Err(); err != nil {
			fmt.Println("Error reading from terminal: ", err.Error())
		}
	}
}

func ConReader(client ClientSelf) { // func for getting server messages
	fmt.Println(client.connection.RemoteAddr())
	buf := make([]byte, 1024)
	for {
		response, err := client.connection.Read(buf)
		if err != nil {
			fmt.Println("Error while reading server msg: ", err.Error())
			time.Sleep(1 * time.Second)
			continue
		}
		fmt.Println(string(buf[:response]))
	}

}

func disconnect(client ClientSelf) {
	client.connection.Write([]byte(client.name + " disconnected"))
	client.connection.Close()
	os.Exit(0)
}
