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

// this will hold all the values recieved from client for serialization and deserialization
type Value struct {
	typ   string  // typ determines the data type carried by Value
	str   string  // str holds the value of string recieved from simple STRING
	num   int     // num holds the integer recieved from INTEGER
	bulk  string  // bulk holds the string recieved from bulk string
	array []Value //array holds the values recieved from the array
}

// Reader struct to read the incoming data/parse and store in Value Struct
type Resp struct {
	reader *bufio.Reader
}

// connection buffer will be passed to NewResp to read and extract data
func NewResp(rd io.Reader) *Resp {
	return &Resp{reader: bufio.NewReader(rd)}
}

// readline reads the line from the buffer
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
	return line[:len(line)-2], n, err
}

func (r *Resp) readInteger() (x, n int, err error) {
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
