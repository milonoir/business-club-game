package network

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net"
	"sync/atomic"
	"time"

	"github.com/gobwas/ws"
	"github.com/gobwas/ws/wsutil"
	"golang.org/x/exp/slog"
)

const (
	pingInterval      = 2 * time.Second
	messageBufferSize = 100
	errClosedByRemote = "use of closed network connection"
	errBrokenPipe     = "write: broken pipe"
)

var (
	compiledClientSideCloseMessage = ws.MustCompileFrame(ws.MaskFrame(ws.NewCloseFrame(ws.NewCloseFrameBody(ws.StatusNormalClosure, ""))))
)

// Connection defines the interface for interacting with a network connection.
type Connection interface {
	Close() error
	Send(Message) error
	Inbox() <-chan Message
	IsAlive() bool
	RemoteAddress() net.Addr
}

// connection implements both server and client side connections.
type connection struct {
	conn     net.Conn
	side     ws.State
	done     chan struct{}
	inbox    chan Message
	alive    atomic.Bool
	lastPong time.Time
	l        *slog.Logger
}

// NewServerConnection configures and returns a connection for the server side.
func NewServerConnection(conn net.Conn, l *slog.Logger) Connection {
	c := &connection{
		conn:  conn,
		side:  ws.StateServerSide,
		done:  make(chan struct{}),
		inbox: make(chan Message, messageBufferSize),
		l: l.With(
			"remote_addr", conn.RemoteAddr(),
			"side", "server",
		),
	}
	c.alive.Store(true)

	go c.receive()
	go c.pinger()

	return c
}

// NewClientConnection configures and returns a connection for the client side.
func NewClientConnection(conn net.Conn, l *slog.Logger) Connection {
	c := &connection{
		conn:  conn,
		side:  ws.StateClientSide,
		done:  make(chan struct{}),
		inbox: make(chan Message, messageBufferSize),
		l: l.With(
			"remote_addr", conn.RemoteAddr(),
			"side", "client",
		),
	}
	c.alive.Store(true)

	go c.receive()

	return c
}

// Send sends a Message to the remote connection.
func (c *connection) Send(msg Message) error {
	b, err := json.Marshal(msg)
	if err != nil {
		return fmt.Errorf("corrupt message: %w", err)
	}

	c.l.Debug("write message", "type", msg.Type(), "payload", msg.Payload())
	if err = wsutil.WriteMessage(c.conn, c.side, ws.OpText, b); err != nil {
		return fmt.Errorf("write message to %s: %w", c.conn.RemoteAddr(), err)
	}

	return nil
}

// Close closes the connection along with all running goroutines.
func (c *connection) Close() error {
	close(c.done)
	close(c.inbox)

	if c.alive.Load() {
		m := ws.CompiledCloseNormalClosure
		if c.side == ws.StateClientSide {
			// NOTE: This is a workaround for the client side close message. It is likely a bug in gobwas/ws.
			m = compiledClientSideCloseMessage
		}
		if _, err := c.conn.Write(m); err != nil {
			return fmt.Errorf("send compiled close to %s: %w", c.conn.RemoteAddr(), err)
		}
	}
	if err := c.conn.Close(); err != nil {
		return fmt.Errorf("close connection %s: %w", c.conn.RemoteAddr(), err)
	}

	return nil
}

// Inbox returns the channel of incoming messages.
func (c *connection) Inbox() <-chan Message {
	return c.inbox
}

// IsAlive returns if the wrapped connection is alive.
func (c *connection) IsAlive() bool {
	return c.alive.Load()
}

func (c *connection) RemoteAddress() net.Addr {
	return c.conn.RemoteAddr()
}

// receive is the message receiver goroutine.
func (c *connection) receive() {
	var (
		err   error
		opErr *net.OpError
		msg   Message
		raw   = make([]wsutil.Message, 0, messageBufferSize)
	)

	for {
		select {
		case <-c.done:
			return
		default:
		}

		if !c.alive.Load() {
			return
		}

		raw = raw[:0]
		raw, err = wsutil.ReadMessage(c.conn, c.side, raw)
		switch {
		case err == nil:
		case errors.Is(err, io.EOF), errors.Is(err, io.ErrUnexpectedEOF), errors.Is(err, net.ErrClosed):
			continue
		case errors.As(err, &opErr):
			if opErr.Error() == errClosedByRemote || opErr.Error() == errBrokenPipe {
				// Connection closed by the remote side.
				c.alive.Store(false)
				return
			}
		default:
			c.l.Error("read message", "error", err)
			continue
		}

		for _, rm := range raw {
			c.l.Debug("raw message", "opcode", rm.OpCode, "payload", rm.Payload)

			// Connection closed by the remote side.
			if rm.OpCode == ws.OpClose {
				c.alive.Store(false)
				_ = c.Close()
				return
			}

			// Ping message - only used by client side.
			if rm.OpCode == ws.OpPing {
				if err = wsutil.WriteMessage(c.conn, c.side, ws.OpPong, rm.Payload); err != nil {
					c.l.Error("write pong", "error", err)
				}
				continue
			}
			// Pong message - only used by server side.
			if rm.OpCode == ws.OpPong {
				c.l.Debug("pong")
				c.lastPong = time.Now()
				continue
			}

			// Text message.
			if msg, err = Parse(rm.Payload); err != nil {
				c.l.Error("parse message", "error", err)
				continue
			}
			select {
			case c.inbox <- msg:
			default:
				// Inbox buffer full, drop message.
			}
		}
	}
}

// pinger is a goroutine used by the server to keep pinging clients in the background.
func (c *connection) pinger() {
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
			if !c.lastPong.IsZero() && c.lastPong.Add(pingInterval*5).Before(time.Now()) {
				// After 5 missed pongs, connection is considered to be dead.
				c.l.Error("ping timeout, connection lost")
				c.alive.Store(false)
				return
			}
			_, err := c.conn.Write(ws.CompiledPing)
			var opErr *net.OpError
			switch {
			case err == nil:
			case errors.As(err, &opErr):
				if opErr.Error() == errBrokenPipe {
					c.l.Error("connection lost", "error", err)
					c.alive.Store(false)
					return
				}
			default:
				c.l.Error("write compiled ping", "error", err)
			}
		}
	}
}
