package taskmaster

import (
	"fmt"
	"log"
	"net"
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

// Shared between client and server
var logger *log.Logger

func init() {
	logger = createLogger(logPath)
}

// App is wrapper for application data
type App struct {
	logger   *log.Logger
	signals  chan os.Signal
	done     chan bool
	oldState *terminal.State
	term     *tty.State
	conn     net.Conn
}

// Logger ...
func Logger() *log.Logger {
	return logger
}

// Create new app (taskmaster)
func Create() *App {
	var err error
	app := &App{
		logger:  logger,
		signals: make(chan os.Signal, 1),
		done:    make(chan bool, 1),
	}
	signal.Notify(app.signals, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	app.oldState, err = terminal.MakeRaw(0)
	if err != nil {
		app.logger.Fatal(err)
	}
	app.term = tty.New(4096)

	// autocompletion
	app.term.SetProposer(func(input string) []string {
		var arr []string
		var result []string
		arr = append(arr, "help")
		arr = append(arr, "status")
		arr = append(arr, "reload")
		arr = append(arr, "start")
		arr = append(arr, "stop")
		for _, item := range arr {
			if strings.HasPrefix(item, input) {
				result = append(result, item)
			}
		}
		return result
	})

	// tcp client
	app.conn, err = net.Dial("tcp", "127.0.0.1:4200")
	if err != nil {
		// app.logger.Println(err)
	}
	return app
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
			if app.conn != nil {
				fmt.Fprintf(app.conn, "exit\n")
			}
			app.logger.Println("exit command")
			break
		} else if input != "" {
			if app.conn != nil {
				fmt.Fprintf(app.conn, input+"\n")
			}
			terminal.Restore(0, app.oldState)
			RunCommand(ParseInput(input))
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

// ParseInput ...
func ParseInput(input string) []string {
	// taskmaster.RealTimeExample()
	if len(input) == 0 {
		return nil
	}
	tokens := strings.Split(input, " ")
	return tokens
}

// RunCommand ...
func RunCommand(tokens []string) {
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
	// file, err := os.OpenFile(filePath, os.O_CREATE|os.O_WRONLY, 0644)
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// logger := log.New(file, "["+time.Now().String()[:19]+"]", 0)
	logger := log.New(os.Stdout, "["+time.Now().String()[:19]+"]", 0)
	logger.Println("Logger created")
	return logger
}
