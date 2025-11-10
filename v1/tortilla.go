// Copyright [2025] Manuel Bovo

// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at

//     http://www.apache.org/licenses/LICENSE-2.0

// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
package v1

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/hashicorp/hcl/v2/gohcl"
	"github.com/hashicorp/hcl/v2/hclparse"
	"github.com/hashicorp/hcl/v2/hclwrite"
	"github.com/rs/zerolog"

	vaultFacade "github.com/mbovo/tortilla/v1/vault"
)

type Tortilla struct {
	ctx             context.Context
	logger          zerolog.Logger
	config          TortillaConfig
	vaultConfigFile string
	Cmd             []string
}

func NewTortilla(ctx context.Context, config TortillaConfig, cmd []string) GenericWrapper {
	return &Tortilla{
		ctx:    ctx,
		logger: *zerolog.Ctx(ctx),
		config: config,
		Cmd:    cmd,
	}
}

func (t *Tortilla) Wrap() error {

	t.logger.Info().Msg("Wrapping Tortilla...")

	level := t.config.VaultLogLevel
	if level == "" {
		level = "error"
	}
	cmdString := "vault agent -config " + t.vaultConfigFile + " -log-level=" + level

	t.logger.Debug().Str("cmdString", cmdString).Send()

	args := strings.Split(cmdString, " ")

	t.logger.Trace().Strs("args", args).Send()

	return executor(t.ctx, args[0], args[1:])

}

func (t *Tortilla) Prepare() (e error) {
	t.logger.Info().Msg("Preparing Tortilla...")

	vault_addr := os.Getenv("VAULT_ADDR")
	if vault_addr == "" {
		e = fmt.Errorf("VAULT_ADDR environment variable is not set")
		t.logger.Err(e).Send()
		return
	}

	_, e = exec.LookPath(t.Cmd[0])
	if e != nil {
		t.logger.Err(e).Str("command", t.Cmd[0]).Send()
	}
	_, e = exec.LookPath("vault")
	if e != nil {
		t.logger.Err(e).Str("command", "vault").Send()
	}
	return
}

func (t *Tortilla) Cook() (err error) {

	t.logger.Info().Msg("Cooking Tortilla...")

	f, err := os.CreateTemp("", "vault-agent-config-*.hcl")
	if err != nil {
		return
	}
	f.Close()
	t.vaultConfigFile = f.Name()

	t.logger.Debug().Str("vaultConfigFile", t.vaultConfigFile).Send()

	cmdString := "vault agent generate-config -type env-template"

	args := strings.Split(cmdString, " ")
	args = append(args, "-exec", fmt.Sprintf("%s", strings.Join(t.Cmd, " ")))

	t.logger.Debug().Interface("secrets", t.config.Secrets).Send()

	t.logger.Debug().Strs("args", args).Send()
	for _, secret := range t.config.Secrets {

		args = append(args, "-path", secret.Path)
		t.logger.Trace().Str("path", secret.Path).Send()
	}
	args = append(args, f.Name())

	t.logger.Debug().Strs("final-args", args).Send()

	// generate the config
	err = executor(t.ctx, args[0], args[1:])

	// re-read the config and modify fields
	cfg := &vaultFacade.GeneratedConfig{}
	parser := hclparse.NewParser()

	// e := hclsimple.DecodeFile(t.vaultConfigFile, nil, cfg)
	inFile, _ := parser.ParseHCLFile(t.vaultConfigFile)
	diag := gohcl.DecodeBody(inFile.Body, nil, cfg)
	if diag.HasErrors() {
		t.logger.Error().Str("Errors", diag.Error()).Send()
		for _, d := range diag {
			t.logger.Error().Str("Detail", d.Detail).Send()
		}
		return fmt.Errorf("error reading hcl config file")
	}

	t.logger.Trace().Msg("Modifying Vault Agent config with environment variable mappings")
	for k, template := range cfg.EnvTemplates {
		t.logger.Trace().Str("templateName", template.Name).Msg("Checking template")
		for _, s := range t.config.Secrets {

			if s.EnvName == "" && s.EnvPrefix == "" {
				continue
			}

			t.logger.Trace().Str("secretPath", s.Path[1:]).Msg("Checking secret path")
			if strings.Contains(template.Contents, strings.ReplaceAll(s.Path[1:], "*", "")) {
				if s.EnvName != "" {
					t.logger.Trace().Str("secretPath", s.Path[1:]).Msg("Substituting env name")
					cfg.EnvTemplates[k].Name = s.EnvName
				}
				if s.EnvPrefix != "" {
					t.logger.Trace().Str("secretPath", s.Path[1:]).Msg("Adding env prefix")
					cfg.EnvTemplates[k].Name = s.EnvPrefix + cfg.EnvTemplates[k].Name
				}
			}
		}
	}

	t.logger.Info().Interface("cfg", cfg).Send()

	f, err = os.Create(t.vaultConfigFile)
	if err != nil {
		return
	}
	hclFile := hclwrite.NewEmptyFile()
	gohcl.EncodeIntoBody(cfg, hclFile.Body())
	_, err = hclFile.WriteTo(f) // f opened above as temp file

	return
}

func executor(ctx context.Context, cmd string, cmdArgs []string) error {

	zerolog.Ctx(ctx).Debug().Strs("args", cmdArgs).Str("command", cmd).Msg("exec")
	p := exec.Command(cmd, cmdArgs...)
	p.Stderr = os.Stderr
	p.Stdout = os.Stdout
	p.Stdin = os.Stdin
	p.Env = os.Environ()

	return p.Run()
}
