package analyzer

import "strings"

func init() {
	defaultRegistry.register(TokenizerWhitespace, TokenizerWhitespaceFunc)
}

const (
	TokenizerWhitespace Type = "whitespace"
)

// TokenizerWhitespaceFunc splits string by whitespace characters (see strings.Fields)
func TokenizerWhitespaceFunc(settings map[string]interface{}) (Func, error) {
	return func(s []string) []string {
		if len(s) == 0 {
			return s
		}

		result := make([]string, 0, len(s))

		for _, str := range s {
			result = append(result, strings.Fields(str)...)
		}

		return result
	}, nil
}
