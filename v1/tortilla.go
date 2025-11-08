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

	"github.com/rs/zerolog"
	"github.com/spf13/viper"
)

type Tortilla struct {
	ctx             context.Context
	logger          zerolog.Logger
	config          *viper.Viper
	vaultConfigFile string
	Cmd             []string
}

func NewTortilla(ctx context.Context, config *viper.Viper, cmd []string) GenericWrapper {
	return &Tortilla{
		ctx:    ctx,
		logger: *zerolog.Ctx(ctx),
		config: config,
		Cmd:    cmd,
	}
}

func (t *Tortilla) Wrap() error {

	t.logger.Info().Msg("Wrapping Tortilla...")

	level := t.config.GetString("logLevel")
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

	f, err := os.CreateTemp("", "vault-agent-config-")
	if err != nil {
		return
	}
	f.Close()
	t.vaultConfigFile = f.Name()

	t.logger.Debug().Str("vaultConfigFile", t.vaultConfigFile).Send()

	cmdString := "vault agent generate-config -type env-template"

	args := strings.Split(cmdString, " ")
	args = append(args, "-exec", fmt.Sprintf("%s", strings.Join(t.Cmd, " ")))

	t.logger.Debug().Strs("args", args).Send()
	for _, path := range t.config.GetStringSlice("paths") {
		args = append(args, "-path", path)
		t.logger.Trace().Str("path", path).Send()
	}
	args = append(args, f.Name())

	t.logger.Debug().Strs("final-args", args).Send()

	err = executor(t.ctx, args[0], args[1:])

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
