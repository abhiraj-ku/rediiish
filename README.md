# Simple Redis Implementation in Go

=================================

This project implements a basic Redis-like server in Go. It supports a few core commands including `PING`, `SET`, `GET`, `HSET`, and `HGET`. The implementation is structured to handle commands using a custom value type, and it uses synchronization mechanisms to manage concurrent access.

## Features

---

- **PING**: Returns "PONG" to confirm the server is running.
- **SET**: Stores a key-value pair in memory.
- **GET**: Retrieves the value associated with a given key.
- **HSET**: Stores a key-value pair in a hash map (nested map).
- **HGET**: Retrieves the value associated with a key in a hash map.

## Structure

---

### Value Type

The `Value` struct is used to represent various data types received from clients, supporting serialization and deserialization for the RESP (REdis Serialization Protocol).

### Commands

Commands are handled by specific functions defined in the `Handlers` map:

- `ping(args []Value)`: Responds with "PONG" or the passed argument.
- `set(args []Value)`: Sets a key to a specified value.
- `get(args []Value)`: Retrieves the value for a given key.
- `hset(args []Value)`: Sets a key-value pair in a hash.
- `hget(args []Value)`: Retrieves a value from a hash.

### Synchronization

Read-write mutexes are used to manage concurrent read and write access to the in-memory data structures.

### RESP Protocol

The server reads and writes data according to the RESP protocol, which allows clients to communicate with the server in a standardized format.

## Getting Started

### Prerequisites

Make sure you have Go installed on your machine. You can download it from [golang.org](https://golang.org/dl/).

### Running the Server

1.  Clone the repository:

    bash

    Copy code

    `git clone <repository-url>
cd <repository-directory>`

2.  Build and run the server:

    bash

    Copy code

    `go run main.go`

### Testing Commands

You can test the server using a Redis client or by creating your own client that communicates over the RESP protocol. Here are some example commands you can use:

- **PING**:

  bash

  Copy code

  `*1
$4
PING`

- **SET**:

  bash

  Copy code

  `*3
$3
SET
$4
name
$7
Abhishek`

- **GET**:

  bash

  Copy code

  `*2
$3
GET
$4
name`

- **HSET**:

  bash

  Copy code

  `*4
$4
HSET
$5
users
$2
u1
$8
Abhishek`

- **HGET**:

  bash

  Copy code

  `*3
$4
HGET
$5
users
$2
u1`

## Conclusion

This implementation serves as a basic example of how to create a simple key-value store similar to Redis using Go. You can expand this further by adding more commands, persistence, and additional features.

Feel free to contribute or modify the code to suit your needs!
