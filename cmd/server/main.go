package main

import (
	"flag"
	"log"

	"github.com/tuommii/taskmaster/job"
	"github.com/tuommii/taskmaster/logger"
)

func main() {
	configPath := flag.String("config", "./config.example2.json", "path to config file")
	_ = logger.Get()
	log.Print("Systemlog test")
	flag.Parse()
	s := newServer(*configPath, job.LoadAll(*configPath))
	s.listenSignals()
	s.launchTasks()
	s.listenConnections()
}
