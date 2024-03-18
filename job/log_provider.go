package job

var (
	LogProviders = map[string]NewLogProviderFunc{}
)

type (
	LogProvider interface {
		// WriteLine should be thread safe,
		// as it's launched in a goroutine to prevent blocking of output from workers
		WriteLine(line string, meta LogMeta) error

		// Close is called when the job is complete.
		// This is useful for gracefully closing files or connections
		Close() error
	}

	LogLevel string

	NewLogProviderFunc func(config any) LogProvider

	LogProviderSharedConfig struct {
		Stdout *bool `yaml:"stdout"`
		Stderr *bool `yaml:"stderr"`
		System *bool `yaml:"system"`
	}

	LogMeta struct {
		Level    LogLevel
		WorkerID string
		TaskID   string
	}
)

const (
	Stdout LogLevel = "stdout"
	Stderr LogLevel = "stderr"
	// System LogLevel is when FancyJobRunner logs something, like a task start or retry
	System LogLevel = "system"
)
