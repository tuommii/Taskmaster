package job

import (
	"bufio"
	"fmt"
	"log"
	"os/exec"
	"strings"
)

// Process represents runnable process
type Process struct {
	Name    string
	Command string
	Cmd     *exec.Cmd
}

// Launch ...
func (p *Process) Launch() {
	splitted := strings.Fields(p.Command)
	p.Cmd = exec.Command(splitted[0], splitted[1:]...)
	stdout, err := p.Cmd.StdoutPipe()
	if err != nil {
		log.Fatal(err)
	}
	p.Cmd.Start()
	buf := bufio.NewReader(stdout) // Notice that this is not in a loop
	num := 1
	for {
		line, _, err := buf.ReadLine()
		if err != nil {
			return
		}
		// if num > 3 {
		// 	return
		// }
		num++
		fmt.Println(string(line))
	}

}
