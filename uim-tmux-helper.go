package main

import (
	"bufio"
	"os/exec"
	"regexp"
)

func main() {
	backtick := exec.Command("uim-fep-tick")
	stdout, err := backtick.StdoutPipe()
	if err != nil {
		panic(err)
	}
	if err := backtick.Start(); err != nil {
		panic(err)
	}

	reader := bufio.NewReader(stdout)

	mode_re := regexp.MustCompile("Sk(.)")

	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			panic(err)
		}

		matched := mode_re.FindStringSubmatch(line)

		if len(matched) == 2 {
			line = "[" + matched[1] + "]"
		}

		err = exec.Command("tmux", "set", "status-left", line).Run()
		if err != nil {
			panic(err)
		}
	}
}
