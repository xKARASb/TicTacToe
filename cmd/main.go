package main

import (
	"fmt"
	"github.com/xkarasb/TicTacToe/pkg/input"
	"strconv"
	"strings"

	"github.com/xkarasb/TicTacToe/core/engine"
)

func main() {
	fmt.Println("Привет!\nВыбери число:\n1. Быть хостом\n2. Присоединиться")
	in := input.GetUserInput()
	choice, err := in.InputInt(make(chan struct{}))
	if err != nil {
		panic(err)
	}
	switch choice {
	case 1:
		err := engine.StartGame("localhost", 0)
		if err != nil {
			fmt.Println(err)
		}
	case 2:
		for {
			fmt.Println("Введите адрес хоста:")
			address, err := in.InputString(make(chan struct{}))
			if err != nil {
				panic(err)
			}
			parsedAddress := strings.Split(address, ":")
			host := parsedAddress[0]
			port, _ := strconv.Atoi(parsedAddress[1])
			err = engine.JoinGame(host, port)
			if err != nil {
				fmt.Println(err)
			}
		}
	}
}
