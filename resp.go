package main

import (
	"bufio"
	"fmt"
	"strings"
)

func testRead(input string) {
	if input == "" {
		input = "$6\r\nfoobar\r\n"
	}
	reader := bufio.NewReaderSize(strings.NewReader(input), 1024)
	cmdByte, _ := reader.ReadByte()

	var strLen byte

	if cmdByte == '$' {
		strLen, _ = reader.ReadByte()
	}

	reader.Discard(2)
	buffer := make([]byte, strLen)
	reader.Read(buffer)
	fmt.Println(string(buffer))
	reader.Discard(2)
}
