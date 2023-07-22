package server

import (
	"context"
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/gobwas/ws"
	"golang.org/x/exp/slog"
)

// Server is responsible for handling WS connections and passing them over to the lobby.
type Server struct {
	port uint16
	wg   sync.WaitGroup
	srv  *http.Server
	l    *slog.Logger

	lobby *lobby
}

func NewServer(port uint16, l *slog.Logger) *Server {
	return &Server{
		port: port,
		l:    l.With("component", "server"),
	}
}

func (s *Server) handler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		conn, _, _, err := ws.UpgradeHTTP(r, w)
		if err != nil {
			s.l.Error("upgrade HTTP connection", "error", err)
			return
		}

		s.l.Info("new connection", "remote_addr", conn.RemoteAddr())
		s.lobby.joinPlayer(conn)
	}
}

func (s *Server) Start() {
	// Start the lobby.
	s.l.Info("starting lobby")
	s.lobby = newLobby(s.l)
	s.wg.Add(1)
	go func() {
		defer s.wg.Done()
		s.lobby.start()
	}()
	s.l.Info("lobby started")

	// Start the HTTP(S) server.
	s.l.Info("starting server", "port", s.port)
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
			s.l.Error("server error", "error", err)
			return
		}
	}()
	s.l.Info("server started")
}

func (s *Server) Stop() error {
	// Stop the HTTP(S) server.
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	s.l.Info("stopping server")
	if err := s.srv.Shutdown(ctx); err != nil {
		return fmt.Errorf("server shutdown failed: %w", err)
	}
	s.l.Info("server stopped")

	// Stop the lobby.
	s.l.Info("stopping lobby")
	s.lobby.stop()
	s.l.Info("lobby stopped")

	// Wait for all goroutines to return.
	s.wg.Wait()

	return nil
}
