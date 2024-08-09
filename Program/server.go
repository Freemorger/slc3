package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strconv"
	"strings"
)

var clients []TCPClient

type TCPClient struct {
	connection net.Conn
	addr       net.Addr
	name       string
}

func main() {
	// Server info
	const NAME = "localhost"
	const PORT = 8081
	const CONN_TYPE = "tcp"
	var domain = (NAME + ":" + strconv.Itoa(PORT))
	var newClient TCPClient

	// Start server
	listener, err := net.Listen(CONN_TYPE, domain)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer listener.Close()
	fmt.Println(CONN_TYPE, " Listener started on ", domain)

	// Loop for accepting new clients
	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Error while client tried to connect: ", err.Error())
			continue
		}
		defer conn.Close()
		fmt.Println("New connection estabilished: ", conn.RemoteAddr().String())
		newClient = newClientTCP(conn)
		clients = append(clients, newClient)

		go handleRequest(clients[len(clients)-1])
		//go readConsole(clients[len(clients) - 1])
	}
}

// Func for handling user requests (including messages as well)
func handleRequest(client TCPClient) {

	scanner := bufio.NewScanner(client.connection)
	for scanner.Scan() {
		clientsMessage := scanner.Text()
		fmt.Println(client.name, ": ", clientsMessage)
		globalSend(clientsMessage)
	}

	if err := scanner.Err(); err != nil {
		fmt.Println("Error reading: ", err.Error())
	}
}

func readConsole(client TCPClient) {
	fmt.Print("Server > ")
	consoleScanner := bufio.NewScanner(os.Stdin)
	// Loop for reading user input
	for {
		for consoleScanner.Scan() {
			fmt.Print("Server > ")
			text := consoleScanner.Text()
			if strings.ToLower(text) == "exit" {
				client.connection.Close()
				os.Exit(0)
			}
			_, err := client.connection.Write([]byte("SERVER: " + text + "\n"))
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

func globalSend(msg string) {
	for _, curClient := range clients {
		_, err := curClient.connection.Write([]byte(curClient.name + ": " + msg))
		if err != nil {
			fmt.Println("Error resending: ", err, ". Disconnecting ", curClient.name)
			curClient.connection.Close()
			removeClientTCP(curClient)
			continue
		}
	}
}

func newClientTCP(conn net.Conn) TCPClient {
	client := TCPClient{connection: conn, addr: conn.RemoteAddr()}
	buf := make([]byte, 1024)
	dat, err := conn.Read(buf)
	if err != nil {
		fmt.Println("Getting client name error: ", err.Error())
		client.name = client.addr.String()
	} else {
		client.name = string(buf[:dat])
	}

	return client
}

func removeClientTCP(client TCPClient) {
	for i, v := range clients {
		if v == client {
			clients = append(clients[:i], clients[i+1:]...)
		}
	}
}
