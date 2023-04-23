package analyzer

func init() {
	defaultRegistry.register(Nop, NopFunc)
}

const (
	Nop Type = "nop"
)

// NopFunc Does nothing
func NopFunc(settings map[string]interface{}) (Func, error) {
	return func(s []string) []string {
		return s
	}, nil
}
