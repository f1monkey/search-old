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
	"github.com/invopop/validation"
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
		r.Get("/", indexListHandler(usecase.NewIndexList(storage.All)))
		r.Get("/{index}", indexGetHandler(usecase.NewIndexGet(storage.Get)))
		r.Delete("/{index}", indexDeleteHandler(usecase.NewIndexDelete(logger, storage.Delete)))
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
				writeSimpleError(w, http.StatusBadRequest, "Index already exists")
				return
			}

			var ve validation.Errors
			if errors.As(err, &ve) {
				handleErr(w, newRequestValidationErr(ve))
				return
			}

			handleErr(w, err)
			return
		}

		w.WriteHeader(http.StatusCreated)
	}
}

func indexDeleteHandler(indexDeleter *usecase.IndexDelete) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		name := chi.URLParam(r, "index")

		if err := indexDeleter.Delete(name); err != nil {
			if errors.Is(err, storage.ErrNotFound) {
				writeSimpleError(w, http.StatusNotFound, http.StatusText(http.StatusNotFound))
				return
			}

			handleErr(w, err)
			return
		}

		setContentType(w)
		w.WriteHeader(http.StatusOK)
	}
}

func indexGetHandler(indexGetter *usecase.IndexGetter) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		name := chi.URLParam(r, "index")

		result, err := indexGetter.Get(name)
		if err != nil {
			if errors.Is(err, storage.ErrNotFound) {
				writeSimpleError(w, http.StatusNotFound, http.StatusText(http.StatusNotFound))
				return
			}

			handleErr(w, err)
			return
		}

		data, err := json.Marshal(result)
		if err != nil {
			handleErr(w, errs.Errorf("index marshal err: %w", err))
			return
		}

		setContentType(w)
		w.WriteHeader(http.StatusOK)
		w.Write(data)
	}
}

type IndexListResponse struct {
	Indexes []index.Index `json:"indexes"`
}

func indexListHandler(indexLister *usecase.IndexList) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		result := indexLister.List()

		data, err := json.Marshal(IndexListResponse{Indexes: result})
		if err != nil {
			handleErr(w, errs.Errorf("indexes marshal err: %w", err))
			return
		}

		setContentType(w)
		w.WriteHeader(http.StatusOK)
		w.Write(data)
	}
}
