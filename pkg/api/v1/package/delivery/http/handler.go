package http

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi"
	"go.uber.org/zap"

	"github.com/nikoksr/proji/pkg/api/v1/domain"
	"github.com/nikoksr/proji/pkg/packages"
)

type packageHandler struct {
	logger  *zap.Logger
	manager packages.Manager
}

// Register registers the package handlers to the router.
func Register(logger *zap.Logger, router *chi.Mux, manager packages.Manager) {
	handler := &packageHandler{
		logger:  logger,
		manager: manager,
	}

	router.Get("/api/v1/packages", handler.Fetch)
	router.Get("/api/v1/packages/{label}", handler.GetByLabel)
	router.Post("/api/v1/packages", handler.Create)
	router.Put("/api/v1/packages/{label}", handler.Update)
	router.Delete("/api/v1/packages/{label}", handler.Delete)
}

func (h *packageHandler) Fetch(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	packageList, err := h.manager.Fetch(r.Context())
	if err != nil {
		h.logger.Error("failed to fetch packageList", zap.Error(err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	dataJSON, err := json.Marshal(packageList)
	if err != nil {
		h.logger.Error("failed to marshal packageList", zap.Error(err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	_, _ = w.Write(dataJSON)
}

func (h *packageHandler) GetByLabel(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	label := chi.URLParam(r, "label")

	_package, err := h.manager.GetByLabel(r.Context(), label)
	if err != nil {
		h.logger.Error("failed to get package", zap.Error(err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	dataJSON, err := json.Marshal(_package)
	if err != nil {
		h.logger.Error("failed to marshal package", zap.Error(err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	_, _ = w.Write(dataJSON)
}

func (h *packageHandler) Create(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	_package := domain.PackageAdd{}
	err := json.NewDecoder(r.Body).Decode(&_package)
	if err != nil {
		h.logger.Error("failed to unmarshal package", zap.Error(err))
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err = h.manager.Store(r.Context(), &_package)
	if err != nil {
		h.logger.Error("failed to create package", zap.Error(err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func (h *packageHandler) Update(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	label := chi.URLParam(r, "label")

	_package := domain.PackageUpdate{}
	err := json.NewDecoder(r.Body).Decode(&_package)
	if err != nil {
		h.logger.Error("failed to unmarshal package", zap.Error(err))
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	_package.Label = label

	err = h.manager.Update(r.Context(), &_package)
	if err != nil {
		h.logger.Error("failed to update package", zap.Error(err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (h *packageHandler) Delete(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	label := chi.URLParam(r, "label")

	err := h.manager.Remove(r.Context(), label)
	if err != nil {
		h.logger.Error("failed to delete package", zap.Error(err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
