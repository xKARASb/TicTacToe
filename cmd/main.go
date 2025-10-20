package main

import (
	"fmt"
	"strconv"
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
		parsedAddress := strings.Split(address, ":")
		host := parsedAddress[0]
		port, _ := strconv.Atoi(parsedAddress[1])
		err := engine.JoinGame(host, port)
		if err != nil {
			fmt.Println(err)
		}
	}
}
