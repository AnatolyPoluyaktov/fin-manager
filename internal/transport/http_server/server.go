package http_server

import (
	"context"
	"net/http"
	"time"
)

// Server представляет HTTP-сервер
type Server struct {
	httpServer *http.Server
}

// NewServer создает и настраивает HTTP-сервер
func NewServer(router http.Handler, address string) *Server {
	return &Server{
		httpServer: &http.Server{
			Addr:           address,
			Handler:        router,
			ReadTimeout:    5 * time.Second,  // Таймаут на чтение запроса
			WriteTimeout:   10 * time.Second, // Таймаут на запись ответа
			IdleTimeout:    60 * time.Second, // Таймаут на ожидание следующего запроса
			MaxHeaderBytes: 1 << 20,          // 1 MB
		},
	}
}

// Start запускает HTTP-сервер
func (s *Server) Start() error {
	return s.httpServer.ListenAndServe()
}

// Shutdown корректно завершает работу сервера
func (s *Server) Shutdown(ctx context.Context) error {
	return s.httpServer.Shutdown(ctx)
}
