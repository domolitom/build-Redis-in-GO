package main

import (
	"bufio"
	"io"
	"strconv"
)

const (
	STRING  = '+'
	ERROR   = '-'
	INTEGER = ':'
	BULK    = '$'
	ARRAY   = '*'
)

type Value struct {
	typ   string
	str   string
	num   int
	bulk  string
	array []Value
}

type Resp struct {
	reader *bufio.Reader
}

func NewResp(reader io.Reader) *Resp {
	return &Resp{reader: bufio.NewReader(reader)}
}

func (r *Resp) readLine() (line []byte, n int, err error) {
	for {
		b, err := r.reader.ReadByte()
		if err != nil {
			return nil, 0, err
		}
		n += 1
		line = append(line, b)
		if len(line) >= 2 && line[len(line)-2] == '\r' {
			break
		}
	}
	return line[:len(line)-2], n, nil
}

func (r *Resp) Read() (Value, error) {

}

func (r *Resp) readInteger() (x int, n int, err error) {
	line, n, err := r.readLine()
	if err != nil {
		return 0, 0, err
	}
	i64, err := strconv.ParseInt(string(line), 10, 64)
	if err != nil {
		return 0, n, err
	}
	return int(i64), n, nil
}

func (r *Resp) readArray() (Value, error) {
	value := Value{}
	value.typ = "array"

	// read array length
	len, _, err := r.readInteger()
	if err != nil {
		return value, err
	}

	// read array values for each line
	value.array = make([]Value, 0)
	for i := 0; i < len; i++ {
		v, err := r.Read()
		if err != nil {
			return value, err
		}
		value.array = append(value.array, v)
	}

	return value, nil
}

func (r *Resp) readBulk() (Value, error) {
	value := Value{}
	value.typ = "bulk"

	//read bulk length
	len, _, err := r.readInteger()
	if err != nil {
		return value, err
	}

	// read bulk value
	bulk := make([]byte, len)

	r.reader.Read(bulk)

	value.bulk = string(bulk)

	r.readLine() // read \r\n

	return value, nil
}
