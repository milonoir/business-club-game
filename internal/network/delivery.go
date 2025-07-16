package network

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net"
	"time"

	"github.com/gobwas/ws"
	"github.com/gobwas/ws/wsutil"
)

type deliveryType string

const (
	deliveryWrapped deliveryType = "wrapped"
	deliveryAck     deliveryType = "ack"

	retryInterval = 2 * time.Second
)

// wrapped adds a unique identifier to the message to allow for tracking delivery and retransmission.
type wrapped struct {
	Id   string
	Type deliveryType
	Msg  json.RawMessage
}

// courier is responsible for managing the delivery of messages and handling retries if necessary.
type courier struct {
	conn net.Conn
	side ws.State
	id   string
	msg  []byte
	done chan struct{}
	l    *slog.Logger
}

// newCourier creates a new courier instance.
func newCourier(conn net.Conn, side ws.State, id string, msg []byte, l *slog.Logger) *courier {
	return &courier{
		conn: conn,
		side: side,
		id:   id,
		msg:  msg,
		done: make(chan struct{}),
		l:    l.With("context", "courier"),
	}
}

// stop stops the courier, closing the done channel to signal it should stop sending messages.
func (c *courier) stop() {
	close(c.done)
}

// run starts the courier, which will attempt to send the message at regular intervals until stopped.
func (c *courier) run() {
	t := time.NewTicker(retryInterval)
	defer t.Stop()

	c.send()
	for {
		select {
		case <-c.done:
			return
		case <-t.C:
			c.send()
		}
	}
}

// send sends the message over the WebSocket connection, logging any errors that occur.
func (c *courier) send() {
	if err := wsutil.WriteMessage(c.conn, c.side, ws.OpText, c.msg); err != nil {
		c.l.Error(fmt.Sprintf("write message %s to %s: %v", c.id, c.conn.RemoteAddr(), err))
	}
}
