package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strconv"
	"strings"
)

type ClientConnection struct {
	Conn      net.Conn
	TCPReader bufio.Reader
}

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

	cliConn := ClientConnection{Conn: conn, TCPReader: *bufio.NewReader(conn)}

	fmt.Println("connection started...")

	command := askCommand(cliConn, []string{"Login", "Create Account"})

	if command == "Create Account" {
		fmt.Println("starting create account...")
		name, err := cliConn.askText("Enter a new Username: ", 3)

		if err != nil {
			fmt.Println(err)
		} else {
			fmt.Println(name)
		}

	} else {
		fmt.Println("closing connection...")
	}
}

func (cliConn ClientConnection) askText(prompt string, minChars int) (string, error) {
	cliConn.Write(prompt)

	for {
		data, err := readClientMessage(cliConn.TCPReader)
		if err != nil {
			return "", err
		}

		if len(data) > minChars {
			return data, nil
		} else {
			cliConn.Write(fmt.Sprintf("must be more than %d characters long\n\n", minChars) + prompt)
		}
	}
}

func askCommand(cliConn ClientConnection, cmds []string) string {
	invalidResponse := "invalid command, try again...\n"

	for {
		var prompt string
		for i := range cmds {
			prompt += fmt.Sprintf("[%d] %s\n", i, cmds[i])
		}
		prompt += "$"

		cliConn.Write(prompt)

		data, err := readClientMessage(cliConn.TCPReader)
		if err != nil {
			fmt.Println(err)
			return ""
		}

		for i := range cmds {
			if isCommandMatch(cmds[i], i, data) {
				return cmds[i]
			}
		}

		fmt.Println("CMD: ", data)
		fmt.Println("****")
		cliConn.Write(invalidResponse)
	}
}

func isCommandMatch(cmd string, index int, input string) bool {
	cmd = strings.ToLower(cmd)
	input = strings.ToLower(input)
	input = strings.TrimSpace(input)

	if strings.Index(cmd, input) == 0 {
		return true
	}

	inputIndex, err := strconv.Atoi(input)
	if err != nil {
		fmt.Println(input, " not a number?")
		return false
	}
	if index == inputIndex {
		return true
	}

	return false
}

func readClientMessage(tcpReader bufio.Reader) (string, error) {
	data, err := tcpReader.ReadString('$')

	if err != nil {
		return "", err
	}

	data = strings.Trim(data, "$")

	fmt.Println("---------")
	fmt.Println("CLIENT")
	fmt.Println("")
	fmt.Print(data)
	fmt.Println("---------")
	return data, nil
}

func (cliConn *ClientConnection) Write(text string) error {
	_, err := cliConn.Conn.Write([]byte(text + "$"))
	return err
}
