package tcp

import (
	"fmt"
	"log"
	"net"

	"github.com/pkg/errors"
)

type Server struct {
	port     string
	tcpAddr  *net.TCPAddr
	listener *net.TCPListener
}

func NewServer(port string) (Server, error) {
	server := Server{port: port}

	err := server.openListener()
	return server, err
}

func (s *Server) openListener() error {
	if s.port == "" {
		log.Fatal("The Port should not be empty")
	}

	laddr, err := net.ResolveTCPAddr("tcp", s.port)
	if err != nil {
		return errors.Wrap(err, "Error finding TCP port address: "+s.port)
	}

	fd, err := net.ListenTCP("tcp", laddr)
	if err != nil {
		return err
	}
	s.tcpAddr = laddr
	s.listener = fd
	return nil
}

func (s *Server) HandleServerLoop(handler func(conn *net.TCPConn)) {
	defer s.listener.Close()

	fmt.Printf("Listening at (%s)\n", s.tcpAddr)

	for {
		conn, err := s.listener.AcceptTCP()
		if err != nil {
			fmt.Printf("Could not accept connection: %v\n", err)
			conn.Close()
			continue
		}

		go handler(conn)
	}
}
