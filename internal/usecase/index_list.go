package usecase

import (
	"github.com/f1monkey/search/internal/index"
)

type IndexList struct {
	lister indexLister
}

type indexLister func() []index.Index

func NewIndexList(lister indexLister) *IndexList {
	return &IndexList{
		lister: lister,
	}
}

func (u *IndexList) List() []index.Index {
	return u.lister()
}
