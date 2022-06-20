package main

import (
	"log"
	"net"

	"github.com/gobwas/ws"
	"github.com/gobwas/ws/wsutil"
	"github.com/milonoir/bc-server/message"
)

type connection struct {
	conn net.Conn
	done chan struct{}

	// channels for communicating with the server/lobby.
	in  <-chan message.Message
	out chan<- message.Message
}

func newConnection() *connection {
	return &connection{}
}

func (c *connection) ping() {
	for {
		select {
		case <-c.done:
			return
		default:
			if _, err := c.conn.Write(ws.CompiledPing); err != nil {
				log.Printf("ERROR - sending ping to client (%s): %v", c.conn.RemoteAddr(), err)
			}
		}
	}
}

func (c *connection) reader() {
	for {
		data, op, err := wsutil.ReadClientData(c.conn)
		if err != nil {
			log.Printf("ERROR - read client (%s) data: %v", c.conn.RemoteAddr(), err)
		}

		_ = data

		switch op {
		case ws.OpPong:
		case ws.OpClose:
		case ws.OpText:
		default:

		}
	}
}
