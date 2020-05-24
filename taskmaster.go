package taskmaster

import (
	"log"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/tuommii/taskmaster/cli"
	"github.com/tuommii/taskmaster/tty"
	"golang.org/x/crypto/ssh/terminal"
)

const logPath = "/tmp/taskmaster_log"

// App is wrapper for application data
type App struct {
	logger   *log.Logger
	signals  chan os.Signal
	done     chan bool
	oldState *terminal.State
	term     *tty.State
}

// Create new app (taskmaster)
func Create() *App {
	var err error
	app := &App{
		logger:  createLogger(logPath),
		signals: make(chan os.Signal, 1),
		done:    make(chan bool, 1),
	}
	signal.Notify(app.signals, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	app.oldState, err = terminal.MakeRaw(0)
	if err != nil {
		app.logger.Fatal(err)
	}
	app.term = tty.New(4096)
	return app
}

func parseInput(input string) []string {
	// taskmaster.RealTimeExample()
	if len(input) == 0 {
		return nil
	}
	tokens := strings.Split(input, " ")
	return tokens
}

func runCommand(tokens []string) {
	if len(tokens) == 0 {
		return
	}
	for _, cmd := range cli.Commands {
		if tokens[0] == cmd.Name || tokens[0] == cmd.Alias {
			cmd.Run(cmd, tokens[1:])
		}
	}
}

func createLogger(filePath string) *log.Logger {
	file, err := os.OpenFile(filePath, os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal(err)
	}
	logger := log.New(file, time.Now().String()[:19]+" ", 0)
	logger.Println("Logger created")
	return logger
}

// ListenSignals ...
func (app *App) ListenSignals() {
	sig := <-app.signals
	app.logger.Println("RECEIVED:", sig)
	app.done <- true
}

// ReadInput ...
func (app *App) ReadInput() {
	for {
		input := app.term.ReadKey(app.signals)
		if input == "exit" {
			app.logger.Println("exit command")
			break
		}
		if input != "" {
			terminal.Restore(0, app.oldState)
			runCommand(parseInput(input))
		}
		terminal.MakeRaw(0)
	}
	app.done <- true
}

// Quit ...
func (app *App) Quit() {
	<-app.done
	terminal.Restore(0, app.oldState)
	app.logger.Println("quit")
	os.Exit(1)
}
