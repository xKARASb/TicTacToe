package server

import (
	"bufio"
	"fmt"
	"net"
)

type GameServer struct {
	host     string
	port     int
	conn     net.Conn
	listener net.Listener
}

func NewGameServer(host string, port int) *GameServer {
	return &GameServer{
		host, port, nil, nil,
	}
}

func (s *GameServer) StartServer(cmd chan string) error {
	tcpAddr, err := net.ResolveTCPAddr("tcp4", fmt.Sprintf("%s:%d", s.host, s.port))

	if err != nil {
		return err
	}

	listener, err := net.ListenTCP("tcp", tcpAddr)
	if err != nil {
		return err
	}
	for {

		conn, err := listener.Accept()
		if err != nil {
			fmt.Println(err)
		}
		handleConnection(conn, cmd)
	}
}

func handleConnection(conn net.Conn, ch chan string) {
	defer conn.Close()

	for {
		data, err := bufio.NewReader(conn).ReadString('\n')
		if err != nil {
			fmt.Println(err)
			return
		}
		ch <- data
	}
}

func (s *GameServer) GetAddr() string {
	return s.listener.Addr().String()
}

func (s *GameServer) Send(msg string) error {
	if _, err := s.conn.Write([]byte(msg)); err != nil {
		return err
	}
	return nil
}
