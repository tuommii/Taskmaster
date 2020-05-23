package job

import (
	"bufio"
	"fmt"
	"log"
	"os/exec"
)

func Testi() {

	dateCmd := exec.Command("ls")
	dateOut, err := dateCmd.Output()
	if err != nil {
		fmt.Println(err)
	}
	str := string(dateOut)
	fmt.Printf("\n%s", str)
	// fmt.Printf("\n%s", strings.TrimSuffix(str, "\n"))

}

func RealTimeExample() {
	cmd := exec.Command("ping", "127.0.0.1")
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		log.Fatal(err)
	}
	cmd.Start()

	buf := bufio.NewReader(stdout) // Notice that this is not in a loop
	num := 1
	for {
		line, _, _ := buf.ReadLine()
		if num > 3 {
			return
		}
		num++
		fmt.Println(string(line))
	}
}
