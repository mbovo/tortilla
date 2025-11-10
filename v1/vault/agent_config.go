package vault

type GeneratedConfig struct {
	AutoAuth       GeneratedConfigAutoAuth       `hcl:"auto_auth,block"`
	TemplateConfig GeneratedConfigTemplateConfig `hcl:"template_config,block"`
	Vault          GeneratedConfigVault          `hcl:"vault,block"`
	EnvTemplates   []GeneratedConfigEnvTemplate  `hcl:"env_template,block"`
	Exec           GeneratedConfigExec           `hcl:"exec,block"`
}

type GeneratedConfigTemplateConfig struct {
	StaticSecretRenderInterval string `hcl:"static_secret_render_interval"`
	ExitOnRetryFailure         bool   `hcl:"exit_on_retry_failure"`
	MaxConnectionsPerHost      int    `hcl:"max_connections_per_host"`
}

type GeneratedConfigExec struct {
	Command                []string `hcl:"command"`
	RestartOnSecretChanges string   `hcl:"restart_on_secret_changes"`
	RestartStopSignal      string   `hcl:"restart_stop_signal"`
}

type GeneratedConfigEnvTemplate struct {
	Name              string `hcl:"name,label"`
	Contents          string `hcl:"contents,attr"`
	ErrorOnMissingKey bool   `hcl:"error_on_missing_key"`
}

type GeneratedConfigVault struct {
	Address string `hcl:"address"`
}

type GeneratedConfigAutoAuth struct {
	Method GeneratedConfigAutoAuthMethod `hcl:"method,block"`
}

type GeneratedConfigAutoAuthMethod struct {
	Type   string                              `hcl:"type"`
	Config GeneratedConfigAutoAuthMethodConfig `hcl:"config,block"`
}

type GeneratedConfigAutoAuthMethodConfig struct {
	TokenFilePath string `hcl:"token_file_path"`
}
