package main

import (
	"flag"

	"github.com/tuommii/taskmaster/job"
)

func main() {
	configPath := flag.String("config", "./config.example.json", "path to config file")
	flag.Parse()
	s := newServer(*configPath, job.LoadAll(*configPath))
	s.listenSignals()
	s.launchTasks()
	s.listenConnections()
}
