package main

import (
	"log"

	"github.com/tuommii/taskmaster"
	"github.com/tuommii/taskmaster/logger"
)

func main() {
	app := taskmaster.Create()
	_ = logger.Get()
	log.Println("CLIENT TEST")
	app.AddLoggerPrefix("CLIENT: ")
	go app.ListenSignals()
	go app.ReadInput()
	app.Quit()
}
