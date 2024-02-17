package healthz

import (
	"log/slog"
	"net/http"
	"sync"
)

type Check func() error

type Handler interface {
	SetLivenessCheck(Check)
	SetReadinessCheck(Check)

	LiveEndpoint(http.ResponseWriter, *http.Request)
	ReadyEndpoint(http.ResponseWriter, *http.Request)
}

type handler struct {
	http.ServeMux

	mutex sync.RWMutex

	livenessCheck  Check
	readinessCheck Check
}

func NewHandler() Handler {
	h := &handler{}

	h.HandleFunc("/live", h.LiveEndpoint)
	h.HandleFunc("/ready", h.ReadyEndpoint)

	return h
}

func (h *handler) SetLivenessCheck(c Check) {
	h.mutex.Lock()
	defer h.mutex.Unlock()
	h.livenessCheck = c
}

func (h *handler) SetReadinessCheck(c Check) {
	h.mutex.Lock()
	defer h.mutex.Unlock()
	h.readinessCheck = c
}

func (h *handler) LiveEndpoint(w http.ResponseWriter, r *http.Request) {
	h.mutex.RLock()
	defer h.mutex.RUnlock()
	h.handle(w, r, h.livenessCheck)
}

func (h *handler) ReadyEndpoint(w http.ResponseWriter, r *http.Request) {
	h.mutex.RLock()
	defer h.mutex.RUnlock()
	h.handle(w, r, h.readinessCheck)
}

func (h *handler) handle(w http.ResponseWriter, r *http.Request, check Check) {
	if r.Method != http.MethodGet {
		status := http.StatusMethodNotAllowed
		http.Error(w, http.StatusText(status), status)
		return
	}

	var err error
	if check != nil {
		err = check()
	}

	var status int = http.StatusOK
	var message string = http.StatusText(status)

	if err != nil {
		status = http.StatusServiceUnavailable
		message = err.Error()
	}

	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(status)
	_, err = w.Write([]byte(message))
	if err != nil {
		slog.Error("Failed to answer health check", "error", err)
	}
}
