package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"sync"
	"time"

	"github.com/gobwas/ws"
)

type server struct {
	port uint16
	wg   sync.WaitGroup

	*lobby
}

func newServer(port uint16) *server {
	return &server{
		port:  port,
		lobby: newLobby(),
	}
}

func (s *server) handler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		conn, _, _, err := ws.UpgradeHTTP(r, w)
		if err != nil {
			log.Printf("upgrade HTTP conn error: %+v", err)
			return
		}

		s.lobby.joinPlayer(conn)
	}
}

func (s *server) start(sig <-chan os.Signal) {
	// Start the lobby.
	log.Println("starting lobby")
	s.wg.Add(1)
	go func() {
		defer s.wg.Done()
		s.lobby.start()
	}()
	log.Println("lobby started")

	// Start the HTTP server.
	log.Printf("starting server at :%d", s.port)
	srv := &http.Server{
		Addr:    fmt.Sprintf(":%d", s.port),
		Handler: s.handler(),
	}
	s.wg.Add(1)
	go func() {
		defer s.wg.Done()
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("server error: %+v", err)
		}
	}()
	log.Println("server started")

	// Wait until server is terminated.
	<-sig

	// Stop the HTTP server.
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	log.Println("stopping server")
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("server shutdown failed: %+v", err)
	}
	log.Println("server stopped")

	// Stop the lobby.
	log.Println("stopping lobby")
	s.lobby.stop()
	log.Println("lobby stopped")

	// Wait for all goroutines to return.
	s.wg.Wait()
}
