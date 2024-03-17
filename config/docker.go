package config

type (
	DockerConfig struct {
		Image       string                `yaml:"image"`
		Registry    *DockerRegistryConfig `yaml:"registry"`
		Environment map[string]string     `yaml:"environment,omitempty"`
	}

	DockerRegistryConfig struct {
		Server   *string `yaml:"server"`
		User     *string `yaml:"user"`
		Password *string `yaml:"password"`
	}
)
