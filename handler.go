package main

import "sync"

var Handlers = map[string]func([]Value) Value{
	"PING": ping,
	"SET":  set,
	"GET":  get,
	"HSET": hset,
	"HGET": hget,
}

var SETs = map[string]string{}
var SETsMutex = &sync.RWMutex{}
var HSETs = map[string]map[string]string{}
var HSETsMutex = &sync.RWMutex{}

func ping(args []Value) Value {
	if len(args) > 1 {
		return Value{typ: "error", str: "ERR wrong number of arguments for 'ping' command"}
	}

	if len(args) == 0 {
		return Value{typ: "string", str: "PONG"}
	}
	return Value{typ: "string", str: args[0].bulk}
}

func set(args []Value) Value {
	if len(args) != 2 {
		return Value{typ: "error", str: "ERR wrong number of arguments for 'set' command"}
	}

	key := args[0].bulk
	value := args[1].bulk

	SETsMutex.Lock()
	SETs[key] = value
	SETsMutex.Unlock()

	return Value{typ: "string", str: "OK"}
}

func get(args []Value) Value {
	if len(args) != 1 {
		return Value{typ: "error", str: "ERR wrong number of arguments for 'get' command"}
	}

	SETsMutex.RLock()
	value, ok := SETs[args[0].bulk]
	SETsMutex.RUnlock()

	if !ok {
		return Value{typ: "nil", str: "(nil)"}
	}

	return Value{typ: "string", str: value}
}

func hset(args []Value) Value {
	if len(args) < 3 {
		return Value{typ: "error", str: "ERR not enough arguments for 'hset' command. Usage: hset id field1 value1 field2 value2 ..."}
	}

	//odd number of arguments needed
	if len(args)%2 == 0 {
		return Value{typ: "error", str: "ERR incorrect nr of arguments for 'hset' command. Usage: hset id field1 value1 field2 value2 ..."}
	}

	id := args[0].bulk
	for i := 1; i < len(args); i += 2 {
		field := args[i].bulk
		value := args[i+1].bulk

		HSETsMutex.Lock()
		if _, ok := HSETs[id]; !ok {
			HSETs[id] = map[string]string{}
		}
		HSETs[id][field] = value
		HSETsMutex.Unlock()
	}

	return Value{typ: "string", str: "OK"}
}

func hget(args []Value) Value {

	if len(args) != 2 {
		return Value{typ: "error", str: "ERR incorrect number of arguments for 'hget', Usage: hget id field"}
	}

	id := args[0].bulk
	field := args[1].bulk
	HSETsMutex.Lock()
	value, ok := HSETs[id][field]
	HSETsMutex.Unlock()

	if !ok {
		return Value{typ: "nil", str: "(nil)"}
	}

	return Value{typ: "string", str: value}

}
