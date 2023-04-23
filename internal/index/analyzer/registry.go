package analyzer

import "github.com/f1monkey/search/pkg/errs"

type registryItem func(map[string]interface{}) (Func, error)

type registry map[Type]registryItem

func (r registry) register(t Type, f registryItem) {
	r[t] = f
}

func (r registry) get(t Type) (registryItem, error) {
	v, ok := r[t]
	if !ok {
		return nil, errs.Errorf("unknown type %q", t)
	}

	return v, nil
}

var defaultRegistry = registry{}
