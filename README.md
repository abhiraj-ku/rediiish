# Simple Redis Implementation in Go

=================================

This project implements a basic Redis-like server in Go. It supports a few core commands including `PING`, `SET`, `GET`, `HSET`, and `HGET`.

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

---

## Future Plans

We plan to extend the functionality of this Redis-like server by adding more commands, including:

`LPUSH`: Inserts one or more values at the head of a list.
`RPUSH`: Inserts one or more values at the tail of a list.
`LPOP` : Removes and returns the first element of a list.
`RPOP` : Removes and returns the last element of a list.
`ZADD` : Adds one or more members to a sorted set, or updates the score of an existing member.

## Contributing

We welcome contributions to this project! Please see our [CONTRIBUTING.md](https://github.com/abhiraj-ku/rediiish/blob/main/CONTRIBUTING.md) for guidelines on how to contribute.

## Code of Conduct

By participating in this project, you agree to abide by our [Code of Conduct](https://github.com/abhiraj-ku/rediiish/blob/main/CODE_OF_CONDUCT.md).
