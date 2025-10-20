package input

import (
	"bufio"
	"io"
	"os"
)

type UserInput struct{}

func (*UserInput) inputWithCancel(cancelChan chan string) string {
	pr, pw := io.Pipe()
	stdChan := make(chan string, 0)

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

func (*UserInput) InputInt() (int, error) {
	return 0, nil
}
