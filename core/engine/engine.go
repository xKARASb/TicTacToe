package engine

import (
	"fmt"
	"math/rand"

	"github.com/xkarasb/TicTacToe/core/client"
	"github.com/xkarasb/TicTacToe/core/game"
	"github.com/xkarasb/TicTacToe/core/server"
	"github.com/xkarasb/TicTacToe/pkg/render"
	"github.com/xkarasb/TicTacToe/pkg/transport"
)

type Player struct {
	mark   string
	buddy  Buddy
	render *render.Window
}

type Buddy interface {
	Send(string) error
}

func StartGame(host string, port int) error {
	srv := server.NewGameServer(host, port)

	clientChan := make(chan string)
	errChan := make(chan error)

	go func() {
		err := srv.StartServer(clientChan)
		if err != nil {
			errChan <- err
			return
		}
	}()
	for !srv.IsListening() {
	}
	fmt.Println(srv.GetAddr())
	go func() {
		err := <-errChan
		if err != nil {
			panic(err)
		}
	}()

	firstMark := "X"
	secondMark := "O"
	if rand.Intn(2) == 1 {
		firstMark = "O"
		secondMark = "X"
	}

	for !srv.IsConnected() {
	}
	err := srv.Send(transport.SetPlayerMsg(secondMark))
	if err != nil {
		return err
	}

	player := &Player{
		mark:   firstMark,
		buddy:  srv,
		render: render.NewWindow(),
	}

	Proccess(clientChan, errChan, player)
	return nil
}

func JoinGame(host string, port string) error {
	clnt := client.NewClient(host, port)

	serverChan := make(chan string)
	errChan := make(chan error)

	go func() {
		err := clnt.Connect(serverChan)
		if err != nil {
			errChan <- err
			return
		}
	}()
	go func() {
		err := <-errChan
		if err != nil {
			panic(err)
		}
	}()

	msg := <-serverChan

	mark, err := transport.ParseSetPlayerMsg(msg)
	if err != nil {
		return err
	}

	player := &Player{
		mark:   mark,
		buddy:  clnt,
		render: render.NewWindow(),
	}

	Proccess(serverChan, errChan, player)
	return nil
}

func Proccess(ch chan string, errCh chan error, player *Player) {
	field := [3][3]string{}
	player.render.Clear()

	turn := false
	if player.mark == "X" {
		turn = true
	}
	player.render.DrawField(field, player.mark)
	if turn {
		UserInput(&turn, &field, player, errCh)
		player.render.Clear()
		player.render.DrawField(field, player.mark)
	}

	for {
		msg := <-ch
		command := transport.ParseCommand(msg)
		switch command {
		case "cell":
			_, err := transport.ParseCellMsg(msg, &field)
			if err != nil {
				errCh <- err
			}

			turn = true
			player.render.Clear()
			player.render.DrawField(field, player.mark)

			res := game.CheckResult(&field)
			switch res {
			case player.mark:
				player.render.Victory()
				err := player.buddy.Send(transport.EndGame)
				if err != nil {
					errCh <- err
				}
			case "draw":
				player.render.Draw()
				err := player.buddy.Send(transport.EndGame)
				if err != nil {
					errCh <- err
				}
			case "ongoing":
				UserInput(&turn, &field, player, errCh)
				player.render.Clear()
				player.render.DrawField(field, player.mark)
			default:
				player.render.Loose()
				err := player.buddy.Send(transport.EndGame)
				if err != nil {
					errCh <- err
				}
			}
		case transport.EndGame:
			res := game.CheckResult(&field)
			switch res {
			case player.mark:
				player.render.Victory()
			case "draw":
				player.render.Draw()
			default:
				player.render.Loose()
			}
		}
	}
}

func UserInput(turn *bool, field *[3][3]string, player *Player, errChan chan error) {
	var (
		x, y int
	)
	for {
		player.render.Turn()
		player.render.InputCoord(true)
		fmt.Scan(&x)
		player.render.InputCoord(false)
		fmt.Scan(&y)
		x--
		y--

		if game.Validate(field, x, y) {
			err := player.buddy.Send(transport.CellMsg(player.mark, x, y))
			if err != nil {
				errChan <- err
			}

			field[x][y] = player.mark
			*turn = false
			break
		} else {
			player.render.IncorrcetInput()
		}
	}
}
