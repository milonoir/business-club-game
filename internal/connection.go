package internal

import (
	"encoding/json"
	"log"
	"net"
	"sync/atomic"
	"time"

	"github.com/gobwas/ws"
	"github.com/gobwas/ws/wsutil"
	message2 "github.com/milonoir/bc-server/internal/message"
)

const (
	pingInterval = 2 * time.Second
)

type connection struct {
	conn net.Conn
	done chan struct{}

	lastPong time.Time
	drop     atomic.Bool
	incoming chan message2.Message
}

func newConnection(conn net.Conn) *connection {
	c := &connection{
		conn:     conn,
		done:     make(chan struct{}),
		incoming: make(chan message2.Message),
	}

	go c.ping()
	go c.receive()

	return c
}

func (c *connection) ping() {
	t := time.NewTicker(pingInterval)
	defer t.Stop()

	for {
		select {
		case <-c.done:
			return
		case <-t.C:
			if _, err := c.conn.Write(ws.CompiledPing); err != nil {
				log.Printf("ERROR - sending ping to client (%s): %v", c.conn.RemoteAddr(), err)
			}
		}
	}
}

func (c *connection) receive() {
	for {
		select {
		case <-c.done:
			return
		default:
		}

		data, op, err := wsutil.ReadClientData(c.conn)
		if err != nil {
			log.Printf("ERROR - read client (%s) data: %v", c.conn.RemoteAddr(), err)
		}

		switch op {
		case ws.OpPong:
			c.lastPong = time.Now()
		case ws.OpText:
			if c.drop.Load() {
				continue
			}
			m, err := message2.Parse(data)
			if err != nil {
				log.Printf("ERROR - parse message (%s): %v", c.conn.RemoteAddr(), err)
				continue
			}
			select {
			case c.incoming <- m:
			default:
			}
		default:
		}
	}
}

func (c *connection) send(msg message2.Message) {
	b, err := json.Marshal(msg)
	if err != nil {
		log.Printf("ERROR - corrupt message: %v", err)
		return
	}

	if err = wsutil.WriteClientMessage(c.conn, ws.OpText, b); err != nil {
		log.Printf("ERROR - write client (%s) message: %v", c.conn.RemoteAddr(), err)
	}
}

func (c *connection) dropMessages() {
	c.drop.Store(true)
}

func (c *connection) keepMessages() {
	c.drop.Store(false)
}
