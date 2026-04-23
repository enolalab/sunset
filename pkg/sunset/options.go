package sunset

// Option configures parsing behavior.
type Option func(*config)

type config struct {
	ForceLanguage string
}

func defaultConfig() *config {
	return &config{}
}

// WithLanguage forces a specific language instead of auto-detecting from extension.
// The id should be a language identifier like "go", "python", "javascript".
func WithLanguage(id string) Option {
	return func(c *config) {
		c.ForceLanguage = id
	}
}
