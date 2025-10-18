package engine

import (
	"bufio"
	"fmt"
	"math/rand"
	"os"

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
	engineExitChan := make(chan struct{})

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

	go func() {
		for {
			fmt.Println("waiting accepting")
			err := srv.AcceptConnection(clientChan)
			fmt.Println("end waiting accepting")
			if err != nil {
				fmt.Println("Accept Error", err)
			}
			fmt.Println("sending to exit chan")
			engineExitChan <- struct{}{}
			fmt.Println("sended")
		}
	}()

	player := &Player{
		mark:   "",
		buddy:  srv,
		render: render.NewWindow(),
	}

	for {
		fmt.Println(srv.IsConnected())
		if !srv.IsConnected() {
			continue
		}

		randomize := rand.Intn(2)
		firstMark, secondMark := game.Marks[randomize], game.Marks[1-randomize]

		err := srv.Send(transport.SetPlayerMsg(secondMark))
		if err != nil {
			return err
		}
		fmt.Println("mark")

		player.mark = firstMark
		fmt.Println("start proccessing")
		Proccess(clientChan, errChan, engineExitChan, player)
		fmt.Println("end proccessing")

		msg := <-clientChan
		fmt.Println("тут", msg)
		switch msg {
		case "restart":
			var restart int
			player.render.RestartRequest()
			fmt.Scan(&restart)
			if restart == 2 {
				return nil
			}
		case transport.Disconnect:
			//👍
			break
		}

	}
}

func JoinGame(host string, port int) error {
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

	player := &Player{
		mark:   "",
		buddy:  clnt,
		render: render.NewWindow(),
	}

	for {
		msg := <-serverChan

		mark, err := transport.ParseSetPlayerMsg(msg)
		if err != nil {
			return err
		}

		player.mark = mark
		Proccess(serverChan, errChan, make(chan struct{}), player)
		var restart int
		player.render.RestartRequest()
		fmt.Scan(&restart)
		if restart == 2 {
			err = player.buddy.Send(transport.Disconnect)
			if err != nil {
				return err
			}
			return nil
		}

		err = player.buddy.Send(transport.Restart)
		if err != nil {
			return err
		}
	}
}

func Proccess(ch chan string, errCh chan error, exitChan chan struct{}, player *Player) {
	field := [3][3]string{}
	player.render.Clear()

	turn := false
	if player.mark == "X" {
		turn = true
	}
	player.render.DrawField(field, player.mark)
	if turn {
		if UserInput(&turn, &field, player, errCh, exitChan) {
			return
		}
		player.render.Clear()
		player.render.DrawField(field, player.mark)
	}

	for {
		select {
		case msg := <-ch:
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
					if UserInput(&turn, &field, player, errCh, exitChan) {
						return
					}
					player.render.Clear()
					player.render.DrawField(field, player.mark)
				default:
					player.render.Loose()
					err := player.buddy.Send(transport.EndGame)
					if err != nil {
						errCh <- err
					}
				}
				if res != "ongoing" {
					return
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
				return
			}
		case <-exitChan:
			return
		}

	}
}

func UserInput(turn *bool, field *[3][3]string, player *Player, errChan chan error, exitChan chan struct{}) bool {
	inputChan := make(chan [2]int)

	f := func(inputChan chan [2]int) {
		reader := bufio.NewReader(os.Stdin)
		var x, y int
		player.render.Turn()
		player.render.InputCoord(true)
		fmt.Fscan(reader, &x)
		player.render.InputCoord(false)
		fmt.Fscan(reader, &y)
		x--
		y--
		inputChan <- [2]int{x, y}

	}

	go f(inputChan)

	for {
		select {
		case <-exitChan:
			return true
		case coords := <-inputChan:
			if game.Validate(field, coords[0], coords[1]) {
				err := player.buddy.Send(transport.CellMsg(player.mark, coords[0], coords[1]))
				if err != nil {
					errChan <- err
				}

				field[coords[0]][coords[1]] = player.mark
				*turn = false
				return false
			} else {
				player.render.IncorrcetInput()
				go f(inputChan)
			}
		}

	}
}
