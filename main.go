package main

import (
	"fmt"
	"io"
	"net/http"
	"time"

	"golang.org/x/net/websocket"
)

type Server struct {
	cons map[*websocket.Conn]bool
}

func NewServer() *Server {
	return &Server{
		cons: make(map[*websocket.Conn]bool),
	}
}

func (s *Server) handleWSOrderBook(ws *websocket.Conn) {
	fmt.Println("New Conn: ", ws.RemoteAddr())

	for {
		payload := fmt.Sprintf("Order book -> %d\n", time.Now().UnixNano())
		ws.Write([]byte(payload))
		time.Sleep(time.Second * 2)
	}
}

func (s *Server) handleWS(con *websocket.Conn) {
	fmt.Println("Remote Addr", con.RemoteAddr())
	s.cons[con] = true
	s.reedLoop(con)
}

func (s *Server) reedLoop(con *websocket.Conn) {
	buf := make([]byte, 1024)
	for {
		n, err := con.Read(buf)
		if err != nil {
			if err == io.EOF {
				break
			}
			fmt.Println("Error: ", err)
			continue
		}
		msg := buf[:n]
		s.BroadCast(msg)
		fmt.Println(string(msg))
	}
}

func (s *Server) BroadCast(b []byte) {
	for ws := range s.cons {
		go func(ws *websocket.Conn) {
			if _, err := ws.Write(b); err != nil {
				fmt.Println("write error: ", err)
			}
		}(ws)
	}
}
func main() {
	server := NewServer()
	http.Handle("/ws", websocket.Handler(server.handleWS))
	http.Handle("/of", websocket.Handler(server.handleWSOrderBook))
	http.ListenAndServe(":3000", nil)
}
