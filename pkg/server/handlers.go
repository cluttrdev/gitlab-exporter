package server

import (
	"net/http"

	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func (s *Server) MetricsHandler(w http.ResponseWriter, r *http.Request) {
	promhttp.Handler().ServeHTTP(w, r)
}
