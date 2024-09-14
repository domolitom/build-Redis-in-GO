package main

import (
	"fmt"
	"net"
	"strings"
)

func main() {
	fmt.Println("Listening on port :6379")

	// Create a new server
	l, err := net.Listen("tcp", ":6379")
	if err != nil {
		fmt.Println(err)
		return
	}

	// Listen for connections
	conn, err := l.Accept()
	if err != nil {
		fmt.Println(err)
		return
	}

	defer conn.Close()

	resp := NewResp(conn)

	for {
		value, err := resp.Read()
		if err != nil {
			fmt.Println(err)
			return
		}

		if value.typ != "array" {
			fmt.Println("Incorrect request, expected array")
			continue
		}

		if len(value.array) == 0 {
			fmt.Println("Incorrect request, expected non-empty array")
			continue
		}

		command := strings.ToUpper(value.array[0].bulk)
		args := value.array[1:]

		fmt.Println("Command: ", command, " Args: ", args)

		writer := NewWriter(conn)
		handler, ok := Handlers[command]
		if !ok {
			fmt.Println("Unknown command: ", command)
			writer.Write(Value{typ: "error", str: "ERR unknown command"})
			continue
		}
		result := handler(args)
		writer.Write(result)

	}
}
