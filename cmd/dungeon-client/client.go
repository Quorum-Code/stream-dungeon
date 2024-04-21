package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
)

func main() {
	// if len(os.Args) == 1 {
	// 	fmt.Println("host:port required")
	// 	os.Exit(1)
	// }

	tcpAddr, err := net.ResolveTCPAddr("tcp4", ":1200" /*os.Args[1]*/)

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	conn, err := net.DialTCP("tcp", nil, tcpAddr)

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	handleConnection(conn)
}

func handleConnection(conn net.Conn) {
	_, err := conn.Write([]byte("Hello tcp server!\n"))
	fmt.Println("sending...")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	reader := bufio.NewReader(os.Stdin)

	for {
		data, err := bufio.NewReader(conn).ReadString('\n')
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Print("Server: ", string(data))

		fmt.Print("Client(You): ")
		text, _ := reader.ReadString('\n')
		conn.Write([]byte(text))
	}
}
