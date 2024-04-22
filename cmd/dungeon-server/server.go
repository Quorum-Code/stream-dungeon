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
	if len(os.Args) == 1 {
		fmt.Println("host:port required")
		os.Exit(1)
	}

	tcpAddr, err := net.ResolveTCPAddr("tcp4", os.Args[1])

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

	askConnection(conn, []string{"Login", "Create Account"})

	handleNewConnection(conn)

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

func handleNewConnection(conn net.Conn) {
	// user := ""
	// pass := ""

	for {
		_, err := conn.Write([]byte("[C]reate account, [L]ogin\n"))
		if err != nil {
			fmt.Println(err)
			return
		}

		data, err := bufio.NewReader(conn).ReadString('\n')
		if err != nil {
			fmt.Println(err)
			return
		}

		if len(data) == 0 {
			continue
		}

		if data[0] == 'C' {
			fmt.Println("got create command...")
			break
		} else if data[0] == 'L' {
			fmt.Println("got login command...")
			break
		}
	}
}

func askConnection(conn net.Conn, cmds []string) string {
	invalidResponse := "invalid command, try again...\n"

	for {
		var prompt string
		for i := range cmds {
			prompt += fmt.Sprintf("[%d] %s\n", i, cmds[i])
		}
		prompt += "$"

		conn.Write([]byte(prompt))

		data, err := bufio.NewReader(conn).ReadString('$')
		data = strings.Trim(data, "$")
		if err != nil {
			fmt.Println(err)
			return ""
		}

		fmt.Println(data)
		conn.Write([]byte(invalidResponse))
	}
}

func isCommandMatch(cmd string, index int, input string) bool {
	if strings.Index(cmd, input) == 0 {
		return true
	}

	inputIndex, err := strconv.Atoi(input)
	if err != nil {
		return false
	}
	if index == inputIndex {
		return true
	}

	return false
}

func readClientMessage(tcpReader bufio.Reader) error {
	data, err := tcpReader.ReadString('$')

	if err != nil {
		return err
	}

	data = strings.Trim(data, "$")

	fmt.Println("---------")
	fmt.Println("Client:")
	fmt.Print(data)
	fmt.Println("---------")
	return nil
}
