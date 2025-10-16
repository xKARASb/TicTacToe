package main

import (
	"fmt"
	"strings"

	"github.com/xkarasb/TicTacToe/core/engine"
)

func main() {
	fmt.Println("Привет!\nВыбери число:\n1. Быть хостом\n2. Присоединиться")
	var choice int
	fmt.Scan(&choice)
	switch choice {
	case 1:
		err := engine.StartGame("localhost", 0)
		if err != nil {
			fmt.Println(err)
		}
	case 2:
		fmt.Println("Введите адресс хоста:")
		var address string
		fmt.Scan(&address)
		host, port := strings.Split(address, ":")[0], strings.Split(address, ":")[1]
		err := engine.JoinGame(host, port)
		if err != nil {
			fmt.Println(err)
		}
	}
}
