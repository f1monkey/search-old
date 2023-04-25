package usecase

import (
	"github.com/f1monkey/search/internal/index"
	"github.com/invopop/validation"
	"go.uber.org/zap"
)

type IndexCreate struct {
	logger  *zap.Logger
	creator indexCreator
}

type indexCreator func(name string, index index.Index) error

func NewIndexCreate(logger *zap.Logger, creator indexCreator) *IndexCreate {
	if logger == nil {
		logger = zap.NewNop()
	}

	return &IndexCreate{
		logger:  logger,
		creator: creator,
	}
}

func (u *IndexCreate) Create(index index.Index) error {
	if err := validation.Validate(index); err != nil {
		return err
	}

	if err := u.creator(index.Name, index); err != nil {
		return err
	}

	u.logger.Info("index created", zap.String("index", index.Name))

	return nil
}
