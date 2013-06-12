package main

import (
	"bufio"
	"bytes"
	"net"
	"os"
	"os/exec"
	"regexp"
)

type UimHelper struct {
	socketPath  string
	readChannel chan string
	reader      *bufio.Reader
}

func NewUimHelper() *UimHelper {
	helper := new(UimHelper)
	helper.socketPath = os.Getenv("HOME") + "/.uim.d/socket/uim-helper"
	helper.readChannel = make(chan string)

	con, err := net.Dial("unix", helper.socketPath)
	if nil != err {
		panic(err)
	}

	helper.reader = bufio.NewReader(con)

	return helper
}

func (helper *UimHelper) doRead() {
	var buf bytes.Buffer
	sep := []byte{'\n'}

	for {
		line, err := helper.reader.ReadBytes('\n')
		if nil != err {
			panic(err)
		}
		buf.Write(line)

		if bytes.Equal(line, sep) {
			helper.readChannel <- buf.String()
			break
		}
	}
}

func (helper *UimHelper) ReadEvent() string {
	go helper.doRead()
	return <-helper.readChannel
}

func main() {
	uim := NewUimHelper()
	re := regexp.MustCompile("branch\\s+\\S+\\s+(\\S+)")

	for {
		event := uim.ReadEvent()
		matched := re.FindAllStringSubmatch(event, 3)
		mode := matched[1][1]

		cmd := exec.Command("tmux", "set", "status-left", "["+mode+"]")
		err := cmd.Run()
		if nil != err {
			panic(err)
		}
	}
}
