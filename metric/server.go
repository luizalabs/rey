package metric

import (
	"context"
	"fmt"
	"net/http"

	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type Server struct {
	httpServer *http.Server
}

func NewServer(port string) *Server {
	mux := http.NewServeMux()
	mux.Handle("/metrics", promhttp.Handler())

	return &Server{
		httpServer: &http.Server{
			Addr:    fmt.Sprintf(":%s", port),
			Handler: mux,
		},
	}
}

func (s *Server) Run() error {
	return s.httpServer.ListenAndServe()
}

func (s *Server) Stop(ctx context.Context) {
	s.httpServer.Shutdown(ctx)
}
