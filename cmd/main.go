package main

import (
	"fmt"
	"github.com/xkarasb/TicTacToe/core/engine"
	"strings"
)

func main() {
	fmt.Println("Hello\nChoose number:\n1. Host\n2. Connect")
	var choice int
	fmt.Scan(&choice)
	switch choice {
	case 1:
		err := engine.StartGame("localhost", 0)
		if err != nil {
			fmt.Println(err)
		}
	case 2:
		fmt.Println("Enter address:")
		var address string
		fmt.Scan(&address)
		host, port := strings.Split(address, ":")[0], strings.Split(address, ":")[1]
		err := engine.JoinGame(host, port)
		if err != nil {
			fmt.Println(err)
		}
	}
}
