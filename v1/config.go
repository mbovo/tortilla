package v1

type TortillaConfig struct {
	LogLevel      string   `yaml:"logLevel,omitempty"`
	VaultLogLevel string   `yaml:"vaultLogLevel,omitempty"`
	Secrets       []Secret `yaml:"secrets"`
}

type Secret struct {
	Path      string `yaml:"path"`
	EnvName   string `yaml:"envName,omitempty"`
	EnvPrefix string `yaml:"envPrefix,omitempty"`
}
