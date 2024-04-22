package main

import (
	"bufio"
	"fmt"
	"io"
	"net"
	"os"
	"strings"
)

func main() {
	if len(os.Args) == 1 {
		fmt.Println("host:port required")
		os.Exit(1)
	}

	tcpAddr, err := net.ResolveTCPAddr("tcp4", os.Args[1])
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
	reader := bufio.NewReader(os.Stdin)
	tcpReader := bufio.NewReader(conn)

	for {
		err := readServerMessage(*tcpReader)

		if err == io.EOF {
			fmt.Println("Connection closed by server...")
			return
		}

		if err != nil {
			fmt.Println(err)
			return
		}

		fmt.Print("> ")
		text, _ := reader.ReadString('\n')
		conn.Write([]byte(text + "$"))
	}
}

func readServerMessage(tcpReader bufio.Reader) error {
	data, err := tcpReader.ReadString('$')

	if err != nil {
		return err
	}

	data = strings.Trim(data, "$")

	fmt.Println("---------")
	fmt.Println("SERVER")
	fmt.Println("")
	fmt.Println(data)
	fmt.Println("---------")
	return nil
}
