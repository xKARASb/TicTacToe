package client

import (
	"bufio"
	"fmt"
	"log"
	"net"
)

type Client struct {
	host string
	port string

	conn net.Conn
}

func NewClient(host string, port string) *Client {
	return &Client{
		host: host,
		port: port,
		conn: nil,
	}
}

func (c *Client) Connect(ch chan string) error {
	tcpAddr, err := net.ResolveTCPAddr("tcp", fmt.Sprintf("%s:%s", c.host, c.port))
	if err != nil {
		return err
	}

	c.conn, err = net.DialTCP("tcp", nil, tcpAddr)
	if err != nil {
		return err
	}

	err = receive(c.conn, ch)
	if err != nil {
		return err
	}

	return nil
}

func (c *Client) Send(message string) error {
	_, err := c.conn.Write([]byte(message))
	if err != nil {
		return err
	}

	return nil
}

func receive(conn net.Conn, ch chan string) error {
	for {
		data, err := bufio.NewReader(conn).ReadString('\n')
		if err != nil {
			log.Fatal(err)
		}
		ch <- data
	}
}
