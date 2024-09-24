package main

import (
	"sync"
)

// Handler to handle commands based on type(array,bulk etc..)
var Handlers = map[string]func([]Value) Value{
	"PING": ping,
	"SET":  set,
	"GET":  get,
	"HSET": hset,
}

//COMMANDS FOR REDIS CLIENT

// PING- returns "PONG" just to check everything is working fine
// BELOW IS HOW THE COMMANDS LOOK LIKE WITH ARGS : SET NAME ABHISHEK
//
//Value{
//		 typ: "array",
//	     array: []Value{
//			Value{typ: "bulk", bulk: "SET"},
//			Value{typ: "bulk", bulk: "name"},
//			Value{typ: "bulk", bulk: "Abhishek"},
//	 	},
//	}

// Above example looks like this in code
// command := Value{typ: "bulk", bulk: "SET"}.bulk // "SET"

// args := []Value{
// 	Value{typ: "bulk", bulk: "name"},
// 	Value{typ: "bulk", bulk: "Abhishek"},
// }

// PING - RETURNS "PONG" IF NO ARGS IS PASSED
// ELSE RETURNS THE ARGS PASSED
func ping(args []Value) Value {
	if len(args) == 0 {

		return Value{typ: "string", str: "PONG"}
	}
	return Value{typ: "string", str: args[0].bulk}
}

// SET - set name(key) abhishek(value)
/*
The SET command in Redis is a key-value pair.
You can set a key to a specific value at any time
and retrieve it later using the GET command.
USES HASHMAP OF `map[string]string`
*/

var SETs = map[string]string{}
var setMutx = sync.RWMutex{} // read-write mutex- resouces can be read by many , written by one at a time

func set(args []Value) Value {
	if len(args) != 2 {
		return Value{typ: "error", str: "ERR wrong number of arguments for 'set' command"}
	}
	key := args[0].bulk
	value := args[1].bulk

	setMutx.Lock()
	SETs[key] = value

	return Value{typ: "string", str: "OK"}
}

// GET - gets the value of (key) , otherwise returns nil
// `get name: return Abhishek`

func get(args []Value) Value {
	if len(args) != 1 {
		return Value{typ: "error", str: "ERR wrong number of arguments for 'get' command"}

	}
	key := args[0].bulk
	setMutx.RLock()
	value, ok := SETs[key]
	setMutx.RUnlock()
	if !ok {
		return Value{typ: "null"}
	}

	return Value{typ: "bulk", bulk: value}

}

//HSET - set of key-value pairs
// this is basically hashmap of hashmap
// `map[string]map[string]string`

/* The skeleton of hset behind the scenes
{
	"users": {
		"u1": "Abhishek",
		"u2": "Arko",
	},
	"languages": {
		"p1": "Golang",
		"p2": "Javascript",
	},
}

*/

var HSETs = map[string]map[string]string{}
var hsetMutx = sync.RWMutex{}

func hset(args []Value) Value {
	if len(args) != 3 {
		return Value{typ: "error", str: "ERR wrong number of arguments for 'hset' command"}
	}
	hash := args[0].bulk
	key := args[1].bulk
	value := args[2].bulk

	hsetMutx.RLock()
	if _, ok := HSETs[hash]; !ok {
		HSETs[hash] = map[string]string{}
	}
	HSETs[hash][key] = value
	hsetMutx.RUnlock()

	return Value{typ: "string", str: "OK"}
}

//HGET - gets the value of (key) set (hset)
// hget users u1

func hget(args []Value) Value {
	if len(args) != 2 {
		return Value{typ: "error", str: "ERR wrong number of arguments for 'hget' command"}
	}
	hash := args[0].bulk
	key := args[1].bulk

	hsetMutx.RLock()
	value, ok := HSETs[hash][key]

	if !ok {
		return Value{typ: "null"}
	}

	return Value{typ: "bulk", bulk: value}
}
