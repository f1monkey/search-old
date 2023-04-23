package analyzer

import (
	"regexp"

	"github.com/f1monkey/search/pkg/errs"
)

func init() {
	defaultRegistry.register(TokenizerRegexp, TokenizerRegexpFunc)
}

const (
	TokenizerRegexp Type = "regexp"
)

// TokenizerRegexpFunc splits string by regular expression
func TokenizerRegexpFunc(settings map[string]interface{}) (Func, error) {
	var (
		expression string
		ok         bool
	)
	for k, v := range settings {
		if k != "pattern" {
			return nil, errs.Errorf("key %q is not allowed", k)
		}
		expression, ok = v.(string)
		if !ok {
			return nil, errs.Errorf("%q must be a string value", v)
		}
	}
	if expression == "" {
		return nil, errs.Errorf("%q key must be provided", "pattern")
	}

	exp, err := regexp.Compile(expression)
	if err != nil {
		return nil, err
	}

	return func(s []string) []string {
		if len(s) == 0 {
			return s
		}

		result := make([]string, 0, len(s))

		for _, str := range s {
			for _, part := range exp.Split(str, -1) {
				if part == "" {
					continue
				}
				result = append(result, part)
			}
		}

		return result
	}, nil
}
