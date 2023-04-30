package main

import (
	"Lila-Back/internal/server"
	"log"
)

func main() {

	port := "8080"

	serv, err := server.New(port)
	if err != nil {
		log.Fatal(err)
	}
	serv.Start()

}
