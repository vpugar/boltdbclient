package boltdbclient

const (
	// DefaultBindAddress is the default binding interface if none is specified.
	DefaultDir      = ""
	DefaultFilename = "boltdb.db"
)

type Config struct {
	Dir      string `toml:"dir"`
	Filename string `toml:"filename"`
}

// NewConfig returns a new instance of Config with defaults.
func NewConfig() Config {
	return Config{
		Dir:      DefaultDir,
		Filename: DefaultFilename,
	}
}

// WithDefaults takes the given config and returns a new config with any required
// default values set.
func (c *Config) WithDefaults() *Config {
	d := *c
	if d.Dir == "" {
		d.Dir = DefaultDir
	}
	if d.Filename == "" {
		d.Filename = DefaultFilename
	}
	return &d
}

func (config Config) Validate() error {

	return nil
}
