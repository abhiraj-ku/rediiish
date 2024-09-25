package main

import (
	"bufio"
	"os"
	"sync"
	"time"
)

type Aof struct {
	file *os.File
	rd   *bufio.Reader
	mu   sync.Mutex
}

// create NewAof method to be used in main.go when server starts
func NewAof(path string) (*Aof, error) {
	f, err := os.OpenFile(path, os.O_CREATE|os.O_RDWR, 0666)
	if err != nil {
		return nil, err
	}

	aof := &Aof{
		file: f,
		rd:   bufio.NewReader(f),
	}

	// goroutine to sync AOF to disk every 1 second
	go func() {
		for {
			aof.mu.Lock()
			aof.file.Sync()
			aof.mu.Unlock()
			time.Sleep(time.Second)
		}
	}()

	return aof, nil

}

// Method to close AOF file when server closes
func (aof *Aof) Close() error {
	aof.mu.Lock()
	defer aof.mu.Unlock()

	return aof.file.Close()
}

// Method to write the commands sent from client to AOF file
func (aof *Aof) Write(value Value) error {
	aof.mu.Lock()

	defer aof.mu.Unlock()
	_, err := aof.file.Write(value.Marshal())
	if err != nil {
		return err
	}
	return nil
}
