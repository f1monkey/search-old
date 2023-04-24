package node

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
)

type indexStorage interface{}

func indexHandler(logger *zap.Logger, storage indexStorage) func(chi.Router) {
	return func(r chi.Router) {
		r.Get("/", indexListHandler(logger, storage))
		r.Get("/{index}", indexGetHandler(logger, storage))
		r.Delete("/{index}", indexDeleteHandler(logger, storage))
		r.Put("/{index}", indexCreateHandler(logger, storage))
	}
}

func indexCreateHandler(logger *zap.Logger, storage indexStorage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// @todo
	}
}

func indexDeleteHandler(logger *zap.Logger, storage indexStorage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// @todo
	}
}

func indexGetHandler(logger *zap.Logger, storage indexStorage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// @todo
	}
}

func indexListHandler(logger *zap.Logger, storage indexStorage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// @todo
	}
}
