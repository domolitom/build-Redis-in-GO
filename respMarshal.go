package main

import (
	"strconv"
)

func (value Value) Marshal() []byte {
	switch value.typ {
	case "string":
		return value.marshalString()
	case "bulk":
		return value.marshalBulk()
	case "array":
		return value.marshalArray()
	case "null":
		return value.marshalNull()
	case "error":
		return value.marshalError()
	default:
		return []byte{}
	}

}

func (value Value) marshalString() []byte {
	var bytes []byte
	bytes = append(bytes, STRING)
	bytes = append(bytes, value.str...)
	bytes = append(bytes, '\r', '\n')
	return bytes
}

func (value Value) marshalBulk() []byte {
	var bytes []byte
	bytes = append(bytes, BULK)
	bytes = append(bytes, strconv.Itoa(len(value.bulk))...)
	bytes = append(bytes, '\r', '\n')
	bytes = append(bytes, value.bulk...)
	bytes = append(bytes, '\r', '\n')
	return bytes
}

func (value Value) marshalArray() []byte {
	var bytes []byte
	bytes = append(bytes, ARRAY)
	bytes = append(bytes, strconv.Itoa(len(value.array))...)
	bytes = append(bytes, '\r', '\n')
	for _, v := range value.array {
		bytes = append(bytes, v.Marshal()...)
	}
	return bytes
}

func (value Value) marshalNull() []byte {
	return []byte("$-1\r\n")
}

func (value Value) marshalError() []byte {
	var bytes []byte
	bytes = append(bytes, ERROR)
	bytes = append(bytes, value.str...)
	bytes = append(bytes, '\r', '\n')
	return bytes
}
