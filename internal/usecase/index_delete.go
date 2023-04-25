package usecase

import (
	"go.uber.org/zap"
)

type IndexDelete struct {
	logger  *zap.Logger
	deleter indexDeleter
}

type indexDeleter func(name string) error

func NewIndexDelete(logger *zap.Logger, deleter indexDeleter) *IndexDelete {
	if logger == nil {
		logger = zap.NewNop()
	}

	return &IndexDelete{
		logger:  logger,
		deleter: deleter,
	}
}

func (u *IndexDelete) Delete(name string) error {
	if err := u.deleter(name); err != nil {
		return err
	}

	u.logger.Info("index deleted", zap.String("index", name))

	return nil
}
