package input

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
	"sync"
)

type userInput struct {
	pr *io.PipeReader
	pw *io.PipeWriter
}

var (
	in   *userInput
	once sync.Once
)

func GetUserInput() *userInput {
	once.Do(func() {
		pr, pw := io.Pipe()
		go func() {
			io.Copy(pw, os.Stdin)
		}()
		in = &userInput{
			pr, pw,
		}
	})
	return in
}

func (i *userInput) inputWithCancel(cancelChan chan struct{}) (string, error) {
	stdChan := make(chan string)

	go func() {
		reader := bufio.NewReader(i.pr)
		val, err := reader.ReadString('\n')
		if err != nil {
			stdChan <- ""
			return
		}

		if val == "EXIT\n" {
			stdChan <- ""
			return
		}
		stdChan <- val

	}()
	for {
		select {
		case <-cancelChan:
			i.pw.Write([]byte("EXIT\n"))
			return "", fmt.Errorf("exit input")
		case msg := <-stdChan:
			return strings.TrimSuffix(msg, "\n"), nil
		}
	}

}

func (i *userInput) InputInt(cancel chan struct{}) (int, error) {
	raw_i, err := i.inputWithCancel(cancel)
	if err != nil {
		return 0, err
	}
	return strconv.Atoi(raw_i)
}
func (i *userInput) InputString(cancel chan struct{}) (string, error) {
	raw_i, err := i.inputWithCancel(cancel)
	if err != nil {
		return "", err
	}
	if raw_i == "" {
		return "", fmt.Errorf("empty string")
	}
	return raw_i, nil
}
