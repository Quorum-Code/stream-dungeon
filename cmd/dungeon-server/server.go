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

	listener, err := net.ListenTCP("tcp", tcpAddr)

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	fmt.Println("Server started...")

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println(err)
		}

		go handleConnection(conn)
	}
}

func handleConnection(conn net.Conn) {
	defer conn.Close()

	fmt.Println("connection started...")

	for {
		data, err := bufio.NewReader(conn).ReadString('\n')
		if err != nil {
			fmt.Println(err)
			return
		}

		fmt.Print("> ", string(data))

		_, err = conn.Write([]byte("hello client!\n"))
		if err != nil {
			fmt.Println(err)
			return
		}
	}
}
