package v1

type TortillaConfig struct {
	LogLevel        string                 `yaml:"logLevel,omitempty"`
	VaultLogLevel   string                 `yaml:"vaultLogLevel,omitempty"`
	Secrets         []Secret               `yaml:"secrets"`
	Transformations []TransformationConfig `yaml:"transformations,omitempty"`
}

type Secret struct {
	Path string `yaml:"path"`
}

type TransformationConfig struct {
	Type   string `yaml:"type"`
	Match  string `yaml:"match,omitempty"`
	Change string `yaml:"change,omitempty"`
	Path   string `yaml:"path,omitempty"`
}
