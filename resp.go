package main

import (
	"bufio"
	"fmt"
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

// Parsing or deserialization process
func (r *Resp) Read() (Value, error) {
	_type, err := r.reader.ReadByte()
	if err != nil {
		return Value{}, err
	}
	switch _type {
	case ARRAY:
		return r.readArray()
	case BULK:
		return r.readBulk()
	default:
		fmt.Printf("Unknown type: %v", string(_type))
		return Value{}, err

	}
}

// implement readArray function
func (r *Resp) readArray() (Value, error) {
	v := Value{}
	v.typ = "array"

	// read length of the array
	len, _, err := r.readInteger()
	if err != nil {
		return v, err
	}

	// for each line read and parse the value and store to slice
	v.array = make([]Value, 0)
	for i := 0; i < len; i++ {
		val, err := r.Read()
		if err != nil {
			return v, err
		}
		v.array = append(v.array, val)
	}
	return v, nil

}

// Implement readBulk function to read
func (r *Resp) readBulk() (Value, error) {
	v := Value{}

	v.typ = "bulk"
	len, _, err := r.readInteger()
	if err != nil {
		return v, err
	}
	bulk := make([]byte, len)
	r.reader.Read(bulk)

	v.bulk = string(bulk)

	r.readLine() // read \r\n

	return v, nil
}
