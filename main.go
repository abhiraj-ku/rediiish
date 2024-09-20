package main

import (
	"fmt"
	"net"
)

func main() {
	fmt.Println("Listening on port :8001")

	server, err := net.Listen("tcp", ":8001")
	if err != nil {
		fmt.Println(err)
		return
	}

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
		fmt.Println(value)

		conn.Write([]byte("+OK\r\n"))
	}

}
