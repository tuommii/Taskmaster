package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/sevlyar/go-daemon"
	"github.com/tuommii/taskmaster/job"
)

func main() {
	configPath := flag.String("config", "./assets/config.example3.json", "path to config file")
	daemonFlag := flag.Bool("d", false, "run as a daemon")
	// syslogFlag := flag.Bool("syslog", false, "log to syslog")

	flag.Parse()
	// _ = logger.Get(*syslogFlag)

	// This must be runned in main
	if *daemonFlag {
		fmt.Println("Started as a daemon")
		cntxt := &daemon.Context{
			PidFileName: "sample.pid",
			PidFilePerm: 0644,
			LogFileName: "sample.log",
			LogFilePerm: 0640,
			WorkDir:     "./",
			Umask:       027,
			Args:        []string{"[go-sample-daemon]"},
		}

		d, err := cntxt.Reborn()
		if err != nil {
			log.Fatal("Unable to run: ", err)
		}
		if d != nil {
			return
		}
		defer cntxt.Release()
	}

	s := newServer(*configPath, job.LoadAll(*configPath))
	s.listenSignals()
	s.launchTasks()
	s.listenConnections()
}
