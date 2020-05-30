package main

import (
	"log"

	"github.com/tuommii/taskmaster/logger"
)

func main() {
	client := Create()
	_ = logger.Get()
	log.Println("CLIENT TEST")
	go client.ListenSignals()
	go client.ReadInput()
	client.Quit()
}
