package internal

import (
	"encoding/json"
	"errors"
	"io"
	"log"
	"net"
	"sync/atomic"
	"time"

	"github.com/gobwas/ws"
	"github.com/gobwas/ws/wsutil"
	"github.com/milonoir/business-club-game/internal/message"
)

const (
	pingInterval      = 2 * time.Second
	errClosedByClient = "use of closed network connection"
)

type connection struct {
	conn net.Conn
	done chan struct{}

	lastPong time.Time
	drop     atomic.Bool
	alive    atomic.Bool
	incoming chan message.Message
}

func newConnection(conn net.Conn) *connection {
	c := &connection{
		conn:     conn,
		done:     make(chan struct{}),
		incoming: make(chan message.Message),
	}

	c.alive.Store(true)

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
			if !c.alive.Load() {
				return
			}
			if _, err := c.conn.Write(ws.CompiledPing); err != nil {
				log.Printf("ERR - [%s] sending ping to client: %v", c.conn.RemoteAddr(), err)
			}
		}
	}
}

func (c *connection) receive() {
	var (
		err   error
		opErr *net.OpError
		msg   message.Message
		raw   = make([]wsutil.Message, 0, 10)
	)

	for {
		select {
		case <-c.done:
			return
		default:
		}

		raw = raw[:0]
		raw, err = wsutil.ReadClientMessage(c.conn, raw)
		switch {
		case err == nil:
		case errors.Is(err, io.EOF), errors.Is(err, io.ErrUnexpectedEOF), errors.Is(err, net.ErrClosed):
			continue
		case errors.As(err, &opErr):
			if opErr.Error() == errClosedByClient {
				// Connection closed by the client.
				c.alive.Store(false)
				return
			}
		default:
			log.Printf("ERR - [%s] read client message: %v", c.conn.RemoteAddr(), err)
			continue
		}

		for _, rm := range raw {
			// Pong message.
			if rm.OpCode == ws.OpPong {
				log.Printf("[%s] pong", c.conn.RemoteAddr())
				c.lastPong = time.Now()
				continue
			}

			// Text message.
			if c.drop.Load() {
				continue
			}
			if msg, err = message.Parse(rm.Payload); err != nil {
				log.Printf("ERR - [%s] parse message: %v", c.conn.RemoteAddr(), err)
				continue
			}
			select {
			case c.incoming <- msg:
			default:
				// Drop message.
			}
		}
	}
}

func (c *connection) send(msg message.Message) {
	b, err := json.Marshal(msg)
	if err != nil {
		log.Printf("ERR - corrupt message: %v", err)
		return
	}

	if err = wsutil.WriteClientMessage(c.conn, ws.OpText, b); err != nil {
		log.Printf("ERR - [%s] write client message: %v", c.conn.RemoteAddr(), err)
	}
}

func (c *connection) dropMessages() {
	c.drop.Store(true)
}

func (c *connection) keepMessages() {
	c.drop.Store(false)
}
