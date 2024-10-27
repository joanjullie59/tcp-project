package main

import (
	"log"
)

func main() {
	server := NewServer(":3000")
	log.Fatal(server.Start())
}
