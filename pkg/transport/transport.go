package transport

import (
	"fmt"
	"strconv"
	"strings"
)

const (
	Connect    = "connect"
	Disconnect = "disconnect"
	YouTurn    = "yturn"
	StartGame  = "start"
	EndGame    = "end"
	Restart    = "restart"
)

func SetPlayerMsg(player string) string {
	return fmt.Sprintf("you %s", player)
}
func ParseSetPlayerMsg(msg string) (string, error) {
	args := strings.Split(msg, " ")

	if args[0] != "you" {
		return "", fmt.Errorf("incorrect command %s is not you", args[0])
	}
	if args[1] != "O" && args[1] != "X" {
		return "", fmt.Errorf("player is not X or O")
	}

	return args[1], nil
}

func CellMsg(player string, x, y int) string {
	return fmt.Sprintf("cell %s %d %d", player, x, y)
}

func ParseCellMsg(msg string, field *[3][3]string) (*[3][3]string, error) {
	args := strings.Split(msg, " ")

	if args[0] != "cell" {
		return nil, fmt.Errorf("incorrect command %s is not cell", args[0])
	}

	if len(args) != 4 {
		return nil, fmt.Errorf("not enough args. %d must be 4", len(args))
	}

	if args[1] != "O" && args[1] != "X" {
		return nil, fmt.Errorf("player is not X or O")
	}

	x, err := strconv.Atoi(args[2])

	if err != nil {
		return nil, fmt.Errorf("third arg (x) is not int. Must be from 0, 1, 2")
	}

	if x < 0 || x > 2 {
		return nil, fmt.Errorf("third arg (x) is not 0, 1, 2")
	}

	y, err := strconv.Atoi(args[3])

	if err != nil {
		return nil, fmt.Errorf("fourth arg (y) is not int. Must be from 0, 1, 2")
	}
	if y < 0 || y > 2 {
		return nil, fmt.Errorf("fourth arg (y) is not 0, 1, 2")
	}

	field[x][y] = args[1]
	return field, nil
}

func ParseCommand(msg string) string {
	args := strings.Split(msg, " ")
	return args[0]
}
