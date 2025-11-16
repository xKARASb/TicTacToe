package server

import (
	"bufio"
	"fmt"
	"net"
	"strings"
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

func (s *GameServer) StartServer(engineChan chan string) error {
	tcpAddr, err := net.ResolveTCPAddr("tcp4", fmt.Sprintf("%s:%d", s.host, s.port))

	if err != nil {
		return err
	}

	s.listener, err = net.ListenTCP("tcp", tcpAddr)
	if err != nil {
		return err
	}
	return nil
}

func tcpSignalsPipe(engineChan chan string) (tcpChan chan string) {
	tcpChan = make(chan string)

	go func() {
		for {
			engineChan <- <-tcpChan
		}
	}()

	return
}

func (s *GameServer) AcceptConnection(engineChan chan string) (err error) {
	s.conn, err = s.listener.Accept()
	if err != nil {
		return err
	}
	handleConnection(s.conn, tcpSignalsPipe(engineChan))
	s.conn = nil
	return nil
}

func handleConnection(conn net.Conn, ch chan string) {
	defer conn.Close()

	for {
		data, err := bufio.NewReader(conn).ReadString('\n')
		if err != nil {
			return
		}
		ch <- strings.Replace(data, "\n", "", 1)
	}
}

func (s *GameServer) GetAddr() string {
	return s.listener.Addr().String()
}

func (s *GameServer) IsListening() bool {
	return s.listener != nil
}
func (s *GameServer) IsConnected() bool {
	return s.conn != nil
}

func (s *GameServer) Send(msg string) error {
	if _, err := s.conn.Write([]byte(msg + "\n")); err != nil {
		return err
	}
	return nil
}
