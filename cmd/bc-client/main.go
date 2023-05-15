package main

import (
	"context"
	"crypto/tls"
	"log"

	"github.com/gobwas/ws"
	"github.com/gobwas/ws/wsutil"
)

// THIS IS ONLY A TEST CLIENT.

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	ws.DefaultDialer.TLSConfig = &tls.Config{
		InsecureSkipVerify: true,
	}

	// TLS
	//conn, _, _, err := ws.DefaultDialer.Dial(ctx, "wss://localhost:8585")
	//if err != nil {
	//	log.Fatal(err)
	//}

	// Non-TLS
	conn, _, _, err := ws.DefaultDialer.Dial(ctx, "ws://localhost:8585")
	if err != nil {
		log.Fatal(err)
	}

	if err = wsutil.WriteClientMessage(conn, ws.OpText, []byte("hello")); err != nil {
		log.Fatal(err)
	}
}
