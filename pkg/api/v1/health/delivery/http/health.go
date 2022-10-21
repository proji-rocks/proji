package http

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi"
	"go.uber.org/atomic"
	"go.uber.org/zap"
)

type (
	healthHandler struct {
		logger    *zap.Logger
		isHealthy *atomic.Bool
		isReady   *atomic.Bool
	}

	check struct {
		Status string `json:"check"`
	}
)

// Register registers the health and ready handlers to the router.
func Register(logger *zap.Logger, router *chi.Mux, isHealthy, isReady *atomic.Bool) {
	handler := &healthHandler{
		logger:    logger,
		isHealthy: isHealthy,
		isReady:   isReady,
	}

	router.Get("/api/v1/healthz", handler.Health)
	router.Get("/api/v1/readyz", handler.Ready)
	router.Post("/api/v1/readyz/enable", handler.ReadyEnable)
	router.Post("/api/v1/readyz/disable", handler.ReadyDisable)
}

var statusOK = check{Status: "ok"}

func (h *healthHandler) Health(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if !h.isHealthy.Load() {
		w.WriteHeader(http.StatusServiceUnavailable)
		return
	}

	dataJSON, err := json.Marshal(statusOK)
	if err != nil {
		h.logger.Error("failed to marshal health check", zap.Error(err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	_, _ = w.Write(dataJSON)
}

func (h *healthHandler) Ready(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if !h.isReady.Load() {
		w.WriteHeader(http.StatusServiceUnavailable)
		return
	}

	dataJSON, err := json.Marshal(statusOK)
	if err != nil {
		h.logger.Error("failed to marshal ready check", zap.Error(err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	_, _ = w.Write(dataJSON)
}

func (h *healthHandler) ReadyEnable(w http.ResponseWriter, _ *http.Request) {
	h.isReady.Store(true)
	w.WriteHeader(http.StatusAccepted)
}

func (h *healthHandler) ReadyDisable(w http.ResponseWriter, _ *http.Request) {
	h.isReady.Store(false)
	w.WriteHeader(http.StatusAccepted)
}
