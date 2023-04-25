package node

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"

	"github.com/f1monkey/search/internal/index"
	"github.com/f1monkey/search/internal/storage"
	"github.com/f1monkey/search/internal/usecase"
	"github.com/f1monkey/search/pkg/errs"
	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
)

type indexStorage interface {
	Create(key string, value index.Index) error
	Get(key string) (index.Index, error)
	Delete(key string) error
	All() []index.Index
}

func indexesHandler(logger *zap.Logger, storage indexStorage) func(chi.Router) {
	return func(r chi.Router) {
		r.Get("/", indexListHandler(logger, storage))
		r.Get("/{index}", indexGetHandler(logger, storage))
		r.Delete("/{index}", indexDeleteHandler(logger, storage))
		r.Put("/{index}", indexCreateHandler(usecase.NewIndexCreate(logger, storage.Create)))
	}
}

func indexCreateHandler(indexCreator *usecase.IndexCreate) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		name := chi.URLParam(r, "index")

		body, err := io.ReadAll(r.Body)
		defer r.Body.Close()

		if err != nil {
			handleErr(w, errs.Errorf("body read err: %w", err))
			return
		}

		index := index.Index{}
		if err := json.Unmarshal(body, &index); err != nil {
			handleErr(w, errs.Errorf("body unmarshal err: %w", err))
			return
		}
		index.Name = name

		if err := indexCreator.Create(index); err != nil {
			if errors.Is(err, storage.ErrAlreadyExists) {
				writeSimpleError(w, http.StatusBadRequest, "index already exists")
				return
			}

			// @todo handle validation error

			handleErr(w, err)
			return
		}

		w.WriteHeader(http.StatusCreated)
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
