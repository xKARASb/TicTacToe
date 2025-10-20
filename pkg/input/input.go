package input

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strconv"
)

func inputWithCancel(cancelChan chan struct{}) string {
	pr, pw := io.Pipe()
	stdChan := make(chan string)

	go func() {
		io.Copy(pw, os.Stdin)
	}()

	go func() {
		reader := bufio.NewReader(pr)
		val, err := reader.ReadString('\n')

		if err != nil {
			stdChan <- ""
		}

		if val == "EXIT" {
			stdChan <- ""
		}
		stdChan <- val

	}()

	for {
		select {
		case <-cancelChan:
			pw.Write([]byte("EXIT\n"))
			return ""
		case msg := <-stdChan:
			return msg
		}
	}

}

func InputInt(cancel chan struct{}) (int, error) {
	raw_i := inputWithCancel(cancel)
	return strconv.Atoi(raw_i)
}
func InputString(cancel chan struct{}) (string, error) {
	raw_i := inputWithCancel(cancel)
	if raw_i == "" {
		return "", fmt.Errorf("empty string")
	}
	return raw_i, nil
}
