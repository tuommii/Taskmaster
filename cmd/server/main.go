package main

import (
	"flag"
	"log"

	"github.com/tuommii/taskmaster/job"
	"github.com/tuommii/taskmaster/logger"
)

func main() {
	configPath := flag.String("config", "./config.example2.json", "path to config file")
	// daemonFlag := flag.Bool("d", false, "run as a daemon")
	_ = logger.Get()
	log.Print("Systemlog test")
	flag.Parse()

	// cntxt := &daemon.Context{
	// 	PidFileName: "sample.pid",
	// 	PidFilePerm: 0644,
	// 	LogFileName: "sample.log",
	// 	LogFilePerm: 0640,
	// 	WorkDir:     "./",
	// 	Umask:       027,
	// 	Args:        []string{"[go-daemon sample]"},
	// }

	// d, err := cntxt.Reborn()
	// if err != nil {
	// 	log.Fatal("Unable to run: ", err)
	// }
	// if d != nil {
	// 	return
	// }
	// defer cntxt.Release()

	s := newServer(*configPath, job.LoadAll(*configPath))
	s.listenSignals()
	s.launchTasks()
	s.listenConnections()
}
