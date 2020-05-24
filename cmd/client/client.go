package main

import (
	"github.com/tuommii/taskmaster"
)

func main() {
	app := taskmaster.Create()
	app.AddLoggerPrefix("CLIENT: ")
	go app.ListenSignals()
	go app.ReadInput()
	app.Quit()
}
