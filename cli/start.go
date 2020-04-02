package cli

import (
	"fmt"
	"os/exec"
)

// StartCmd implements help command
var StartCmd = &Command{
	Name:  "start",
	Usage: "Start a job",
	Alias: "run",
	Run:   start,
}

func start(cmd *Command, args []string) {
	fmt.Print("\nSTART!")
	for _, arg := range args {
		fmt.Print("\nstart ", arg)
	}
	runCat()
}

// for test
func runCat() {
	catCmd := exec.Command("cat", "Makefile")
	cat, err := catCmd.Output()
	if err != nil {
		fmt.Println(err)
	}
	str := string(cat)
	fmt.Printf("%s", str)

}
