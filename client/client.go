package main

import (
	"fmt"
	"net"
	"time"

	"github.com/rs/zerolog/log"
)

// Client struct  
type Client struct {
	conn net.Conn
}

func main() {
	client := &Client{}
	// Tray Connectiion
	TryConnection(client)

	counter := 0
	for {
		counter++
		s := fmt.Sprintf("%s, %d", "Gelow", counter)
		if err := client.SendMsg([]byte(s)); err != nil {
			TryConnection(client)
		}
		time.Sleep(time.Second * 1)
	}
}

func TryConnection(client *Client) {
	for {
		err := client.DialConnection()
		if err == nil {
			break
		}
		log.Error().Msg(err.Error())
		time.Sleep(time.Second * 1)
	}
}

// SendMsg method  
func (c *Client) SendMsg(msg []byte) error {
	_, err := c.conn.Write(msg)
	if err != nil {
		return err
	}
	return nil
}

func (c *Client) DialConnection() error {
	conn, err := net.Dial("tcp", ":8080")
	if err != nil {
		return err
	}
	c.conn = conn
	return nil
}
