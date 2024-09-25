package main

import (
	"fmt"
	"net"
	"strings"
)

func main() {
	fmt.Println("Listening on port :8001")

	// Create a new server to listen tcp connection
	server, err := net.Listen("tcp", ":8001")
	if err != nil {
		fmt.Println(err)
		return
	}

	// write each request to AOF file sent by client
	aof, err := NewAof("database.aof")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer aof.Close()

	// Reading the value from AOF file
	aof.Read(func(value Value) {
		command := strings.ToUpper(value.array[0].bulk)

		args := value.array[1:]

		handler, ok := Handlers[command]

		if !ok {
			fmt.Println("Invalid command: ", command)
			return
		}
		handler(args)
	})

	// listen for connection
	conn, err := server.Accept()
	if err != nil {
		fmt.Println(err)
		return
	}
	defer conn.Close()

	// infinite loop to recieve client request
	for {
		resp := NewResp(conn)
		value, err := resp.Read()
		if err != nil {
			fmt.Println(err)
			return
		}
		if value.typ != "array" {
			fmt.Println("Invalid request, expecting array")
			continue
		}
		if len(value.array) == 0 {
			fmt.Println("Invalid request, expected array length > 0")
			continue
		}

		commands := strings.ToUpper(value.array[0].bulk)
		args := value.array[1:]

		writer := NewWriter(conn)
		hanlder, ok := Handlers[commands]
		if !ok {
			fmt.Println("Invalid command: ", commands)
			writer.Write(Value{typ: "string", str: ""})
			continue
		}
		if commands == "SET" || commands == "HGET" {
			aof.Write(value)
		}
		result := hanlder(args)
		writer.Write(result)
		// writer.Write(Value{typ: "string", str: "OK"})
	}

}
