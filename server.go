package main

import (
	"encoding/csv"
	"fmt"
	"log"
	"net"
	"os"
	"time"
)

type Server struct {
	listener   net.Listener
	ln         net.Listener
	listenAddr string
	conns      []net.Conn
}

func NewServer(listenerAddr string) *Server {
	return &Server{
		listenAddr: listenerAddr,
	}
}
func (s *Server) Start() error {
	ln, err := net.Listen("tcp", s.listenAddr)
	if err != nil {
		return err
	}
	defer ln.Close()
	s.listener = ln
	s.acceptloop()
	return nil
}
func (s *Server) acceptloop() {
	fmt.Println("Waiting for connections")
	for {
		conn, err := s.listener.Accept()
		if err != nil {
			fmt.Println("accept error:", err)
			continue
		}
		s.conns = append(s.conns, conn)
		fmt.Println("Accept:", conn.RemoteAddr())
		go s.readloop(conn)
	}
}
func (s *Server) readloop(conn net.Conn) {
	defer conn.Close()
	buf := make([]byte, 2048)
	for {
		n, err := conn.Read(buf)
		if err != nil {
			fmt.Println("read  error:", err)
			continue
		}

		msg := string(buf[:n])
		s.broadcastMsg(msg)
		fmt.Println("Read:", msg)

		/// open file
		// append (remoteaddr, datetime, msg) to csv file
		// save file
		csvfile, err := os.OpenFile("messages.csv", os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0600)
		if err != nil {
			log.Fatalf(" failed to open file,:%s", err)
		} else {
			fmt.Println("data csv file open")
		}
		writer := csv.NewWriter(csvfile)

		remoteAddress := conn.RemoteAddr().String()
		dateTime := time.Now().String()
		row := []string{remoteAddress, dateTime, msg}

		err = writer.Write(row)
		if err != nil {
			fmt.Println("write error:", err)
		}
		writer.Flush()
		csvfile.Close()
	}
}

func (s *Server) broadcastMsg(msg string) {
	for idx, conn := range s.conns {
		msg = fmt.Sprintf("%s said '%s'", conn.RemoteAddr(), msg)
		_, err := conn.Write([]byte(msg))
		if err != nil {
			fmt.Println("broadcast error:", err, "removing")
			s.conns = append(s.conns[:idx], s.conns[idx+1:]...) /// deletes the client
		}
	}
}
