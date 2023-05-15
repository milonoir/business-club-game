package internal

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/gobwas/ws"
)

// Server is responsible for handling WS connections and passing them over to the lobby.
type Server struct {
	port uint16
	wg   sync.WaitGroup
	srv  *http.Server

	*lobby
}

func NewServer(port uint16) *Server {
	return &Server{
		port: port,
	}
}

func (s *Server) handler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		conn, _, _, err := ws.UpgradeHTTP(r, w)
		if err != nil {
			log.Printf("upgrade HTTP conn error: %+v", err)
			return
		}

		s.lobby.joinPlayer(conn)
	}
}

func (s *Server) Start() {
	// Start the lobby.
	log.Println("starting lobby")
	s.lobby = newLobby()
	s.wg.Add(1)
	go func() {
		defer s.wg.Done()
		s.lobby.start()
	}()
	log.Println("lobby started")

	// Start the HTTP(S) server.
	log.Printf("starting server at :%d", s.port)
	s.srv = &http.Server{
		Addr:    fmt.Sprintf(":%d", s.port),
		Handler: s.handler(),
	}
	s.wg.Add(1)
	go func() {
		defer s.wg.Done()

		// TLS
		//if err := srv.ListenAndServeTLS(
		//	filepath.Join("cert", "localhost.crt"),
		//	filepath.Join("cert", "localhost.key"),
		//); err != nil && err != http.ErrServerClosed {
		//	log.Fatalf("server error: %+v", err)
		//}

		// Non-TLS
		if err := s.srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("server error: %+v", err)
		}
	}()
	log.Println("server started")
}

func (s *Server) Stop() error {
	// Stop the HTTP(S) server.
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	log.Println("stopping server")
	if err := s.srv.Shutdown(ctx); err != nil {
		return fmt.Errorf("server shutdown failed: %w", err)
	}
	log.Println("server stopped")

	// Stop the lobby.
	log.Println("stopping lobby")
	s.lobby.stop()
	log.Println("lobby stopped")

	// Wait for all goroutines to return.
	s.wg.Wait()

	return nil
}
