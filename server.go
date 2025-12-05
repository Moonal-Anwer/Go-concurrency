package main

import (
	"bufio"
	"fmt"
	"net"
	"sync"
)

type Client struct {
	id   int
	conn net.Conn
	ch   chan string
}

type Server struct {
	mu      sync.Mutex
	clients map[int]*Client
	nextID  int
}

func NewServer() *Server {
	return &Server{
		clients: make(map[int]*Client),
	}
}

func (s *Server) addClient(conn net.Conn) *Client {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.nextID++
	c := &Client{
		id:   s.nextID,
		conn: conn,
		ch:   make(chan string),
	}
	s.clients[c.id] = c
	return c
}

func (s *Server) removeClient(id int) {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.clients, id)
}

func (s *Server) broadcast(senderID int, msg string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	for id, c := range s.clients {
		if id != senderID {
			c.ch <- msg
		}
	}
}

func (s *Server) handleClient(c *Client) {
	go func() {
		for msg := range c.ch {
			fmt.Fprintln(c.conn, msg)
		}
	}()

	// announce join
	s.broadcast(c.id, fmt.Sprintf("User %d joined", c.id))

	scanner := bufio.NewScanner(c.conn)
	for scanner.Scan() {
		text := scanner.Text()
		s.broadcast(c.id, fmt.Sprintf("User %d: %s", c.id, text))
	}

	// client disconnected
	s.removeClient(c.id)
	c.conn.Close()
}

func main() {
	server := NewServer()

	ln, err := net.Listen("tcp", ":1234")
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	fmt.Println("Server running on port 1234...")

	for {
		conn, _ := ln.Accept()
		client := server.addClient(conn)
		go server.handleClient(client)
	}
}
