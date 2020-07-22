package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/sevlyar/go-daemon"
	"github.com/tuommii/taskmaster/job"
	"github.com/tuommii/taskmaster/logger"
)

func main() {
	configPath := flag.String("c", "./assets/config.example4.json", "path to config file")
	daemonFlag := flag.Bool("d", false, "run as a daemon")
	debugFlag := flag.Bool("debug", false, "log to stdout")
	silenceFlag := flag.Bool("s", false, "no logging")
	// syslogFlag := flag.Bool("syslog", false, "log to syslog")

	flag.Parse()

	if *debugFlag {
		logger.ChangeOutput(os.Stdout)
	}
	if *silenceFlag {
		file, _ := os.OpenFile(os.DevNull, 0, 0)
		logger.ChangeOutput(file)
	}

	// This must be runned in main
	if *daemonFlag {
		fmt.Println("Started as a daemon")
		cntxt := &daemon.Context{
			PidFileName: "sample.pid",
			PidFilePerm: 0644,
			LogFileName: "/tmp/sample.log",
			LogFilePerm: 0640,
			WorkDir:     "./",
			Umask:       027,
			Args:        []string{"[go-sample-daemon]"},
		}

		d, err := cntxt.Reborn()
		if err != nil {
			logger.Fatal("Unable to run: ", err)
		}
		if d != nil {
			return
		}
		defer cntxt.Release()
	}

	config := *configPath
	for _, arg := range os.Args[1:] {
		if arg[0] != '-' {
			config = arg
		}
	}

	s := newServer(config, job.LoadAll(config))
	s.listenSignals()
	s.launchTasks()
	s.listenConnections()
}
