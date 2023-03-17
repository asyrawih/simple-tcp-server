// Server TCP
package main

import (
	"fmt"
	"io"
	"net"
	"os"
	"os/signal"
	"syscall"

	"github.com/rs/zerolog/log"
)

// Server struct  
type Server struct {
	Addr string
	Conn []map[string]*net.Conn
}

// Newserver function  
func Newserver() *Server {
	return &Server{
		Addr: "",
		Conn: []map[string]*net.Conn{},
	}
}

// AddConnection method  
func (s *Server) AddConnection(conn *net.Conn) {
	// Create a new map to accept new connection
	mapConn := make(map[string]*net.Conn)
	mapConn[s.Addr] = conn
	s.Conn = append(s.Conn, mapConn)
}

// Read message From Connection
func (s *Server) Read(conn *net.Conn) {
	for {
		// Read message from connection
		buffer := make([]byte, 1024)
		n, err := (*conn).Read(buffer)
		if err != nil {
			if err == io.EOF {
				log.Error().Err(err).Msg("Connection closed")
				s.Close(conn)
				return
			}
			log.Error().Err(err).Msg("Error reading from connection")
			return
		}
		remoteAddr := (*conn).RemoteAddr().String()
		fmt.Println(remoteAddr, string(buffer[:n]))
		msg := fmt.Sprintf("[%s] Reply :%s", (*conn).RemoteAddr().String(), string(buffer[:n]))
		s.Broadcast(conn, msg)
	}
}

// Broadcast method  
func (s *Server) Broadcast(currentConn *net.Conn, msg string) {
	for _, conn := range s.Conn {
		for _, conn := range conn {
			if conn != currentConn {
				_, err := (*conn).Write([]byte(msg))
				if err != nil {
					return
				}
			}
		}
	}
}

// Close connection current connection
func (s *Server) Close(conn *net.Conn) {
	if err := (*conn).Close(); err != nil {
		log.Err(err).Caller().Msg("Error Close Connection")
	}
}

// Listen method  
func (s *Server) Listen(port string) {
	// Listen to port
	listener, err := net.Listen("tcp", port)
	if err != nil {
		log.Error().Err(err).Msg("Error listening to port")
		return
	}

	defer func() {
		if err := listener.Close(); err != nil {
			log.Err(err).Caller().Msg("Error closing listener")
		}
	}()

	log.Info().Msgf("Listening to port %s", port)

	server := new(Server)
	for {
		c, err := listener.Accept()
		server.Addr = c.RemoteAddr().String()
		if err != nil {
			log.Err(err).Caller().Msg("Error creating listener")
		}
		server.AddConnection(&c)
		go server.Read(&c)
	}
}

func main() {
	server := new(Server)
	server.Listen("0.0.0.0:8080")

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)

	<-quit
	// Todo Do Something when get quit signal
	log.Info().Msg("Shutting down server")
}
