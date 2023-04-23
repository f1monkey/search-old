package analyzer

import "github.com/f1monkey/search/pkg/errs"

type Func func([]string) []string
type Handler func(next Func) Func
type Type string

type Analyzer struct {
	Type     Type                   `json:"type"`
	Settings map[string]interface{} `json:"settings"`
}

func New(t Type, settings map[string]interface{}) Analyzer {
	return Analyzer{
		Type:     t,
		Settings: settings,
	}
}

func (a Analyzer) Validate() error {
	_, err := a.Func()
	return err
}

// Func get analyzer func by name
func (a Analyzer) Func() (Func, error) {
	f, err := defaultRegistry.get(a.Type)
	if err != nil {
		return nil, err
	}

	return f(a.Settings)
}

// Chain build a chain from analyzers by their names
func Chain(items []Analyzer) (Func, error) {
	if len(items) == 0 {
		return nil, errs.Errorf("chain cannot be empty")
	}

	var h Func
	for i := len(items) - 1; i >= 0; i-- {
		f, err := items[i].Func()
		if err != nil {
			return nil, err
		}

		h = handler(f, h)
	}

	return h, nil
}

func handler(current Func, next Func) Func {
	if next == nil {
		return current
	}

	return func(s []string) []string {
		return next(current(s))
	}
}
