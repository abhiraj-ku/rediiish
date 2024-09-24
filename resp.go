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

// Value represents a data structure for holding various types of values received from a client.
// It supports serialization and deserialization of different data types.

type Value struct {
	typ   string  // typ indicates the type of data stored in the Value (e.g., STRING, INTEGER, BULK, ARRAY)
	str   string  // str holds a simple string value (used when typ is STRING)
	num   int     // num holds an integer value (used when typ is INTEGER)
	bulk  string  // bulk holds a bulk string value (used when typ is BULK)
	array []Value // array contains a list of Value elements (used when typ is ARRAY)
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

// Implement the Writing Part
// Writer struct
type Writer struct {
	writer io.Writer
}

func NewWriter(w io.Writer) *Writer {
	return &Writer{writer: w}
}

// all marshal methods
func (v Value) Marshal() []byte {
	switch v.typ {
	case "array":
		return v.marshalArray()
	case "bulk":
		return v.marshalBulk()
	case "string":
		return v.marshalString()
	case "error":
		return v.marshalError()
	case "null":
		return v.marshalNull()
	default:
		return []byte{}
	}
}

// marshal string value
func (v Value) marshalString() []byte {
	var bytes []byte
	bytes = append(bytes, STRING)
	bytes = append(bytes, v.str...)
	bytes = append(bytes, '\r', '\n')

	return bytes

}

// marshal Bulk value
func (v Value) marshalBulk() []byte {
	var bytes []byte
	bytes = append(bytes, BULK)
	bytes = append(bytes, strconv.Itoa(len(v.bulk))...)
	bytes = append(bytes, '\r', '\n')
	bytes = append(bytes, v.bulk...)
	bytes = append(bytes, '\r', '\n')

	return bytes
}

// marshal Array value
func (v Value) marshalArray() []byte {
	var bytes []byte
	len := len(v.array)
	bytes = append(bytes, ARRAY)
	bytes = append(bytes, strconv.Itoa(len)...)
	bytes = append(bytes, '\r', '\n')

	for i := 0; i < len; i++ {
		bytes = append(bytes, v.array[i].Marshal()...)
	}

	return bytes
}

// marshal error values
func (v Value) marshalError() []byte {
	var bytes []byte
	bytes = append(bytes, ERROR)
	bytes = append(bytes, v.str...)
	bytes = append(bytes, '\r', '\n')

	return bytes
}

// marshal nil value
func (v Value) marshalNull() []byte {
	return []byte("$-1\r\n")
}

// Method that takes value and writes as output
func (w *Writer) Write(v Value) error {
	var bytes = v.Marshal()
	_, err := w.writer.Write(bytes)
	if err != nil {
		return err
	}
	return nil
}
