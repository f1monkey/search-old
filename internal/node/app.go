package node

import (
	"context"
	"errors"
	"fmt"
	"net"
	"net/http"
	"os"
	"path"

	"github.com/f1monkey/search/internal/index"
	"github.com/f1monkey/search/internal/storage"
	"github.com/f1monkey/search/pkg/errs"
	"github.com/go-chi/chi/v5"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

type Node struct {
	logger *zap.Logger
	server *http.Server
}

func New(ctx context.Context, logger *zap.Logger) (*Node, error) {
	if logger == nil {
		logger = zap.NewNop()
	}

	storagePath := viper.GetString("storage.path")
	if err := os.MkdirAll(storagePath, 0755); err != nil {
		return nil, errs.Errorf("storage path create err: %w", err)
	}

	indexStorage, err := storage.NewAOFFromPath[string, index.Index](path.Join(storagePath, "indexes.dat"))
	if err != nil {
		return nil, err
	}

	mux := chi.NewMux()
	mux.Route("/indexes", indexesHandler(logger, indexStorage))

	return &Node{
		logger: logger,
		server: &http.Server{
			Addr:    viper.GetString("node.server.address"),
			Handler: mux,
			BaseContext: func(net.Listener) context.Context {
				return ctx
			},
		},
	}, nil
}

func (n *Node) Start(ctx context.Context) error {
	n.logger.Info("node starting")

	go func(ctx context.Context) {
		defer panicHandle(ctx, n.logger)
		n.logger.Sugar().Infof("server listening on %s", n.server.Addr)
		err := n.server.ListenAndServe()
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			panic(err)
		}
	}(ctx)

	return nil
}

func (n *Node) Stop(ctx context.Context) error {
	n.logger.Info("node stoppping...")

	n.logger.Info("http server stoppping...")
	if err := n.server.Shutdown(ctx); err != nil {
		n.logger.Error("server shutdown err", zap.Error(err))
		if err := n.server.Close(); err != nil {
			n.logger.Error("server close err", zap.Error(err))
		}
	}
	n.logger.Info("http server stopped")

	n.logger.Info("node stopped")
	return nil
}

func panicHandle(ctx context.Context, l *zap.Logger) {
	if r := recover(); r != nil {
		err, ok := r.(error)
		if !ok {
			err = fmt.Errorf("%v", r)
		}

		l.Fatal("node panic", zap.Error(err))
	}
}
