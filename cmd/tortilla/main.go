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

package main

import (
	"context"
	"fmt"
	"os"
	"strings"

	v1 "github.com/mbovo/tortilla/v1"
	"github.com/rs/zerolog"
	"go.yaml.in/yaml/v4"
)

const (
	configFileName = "tortilla.yaml"
)

var (
	// These are set by the build script
	version = "dev"
	commit  = "HEAD"
	date    = "unknown"

	commandName  = "tortilla"
	commandShort = "A wrap-per for vault-agent to manage secrets easily."
	commandLong  = `====================================================================
 It's a Wrap!

 Tortilla: A wrap-per for vault-agent to manage secrets easily.
====================================================================
Usage: tortilla <command> [arguments]`

	commandVersion = fmt.Sprintf("%s - %s - %s", version, commit, date)

	ctx = zerolog.New(
		zerolog.ConsoleWriter{
			Out:          os.Stdout,
			PartsExclude: []string{zerolog.LevelFieldName, zerolog.TimestampFieldName},
		},
	).WithContext(context.Background())
)

func init() {
	initConfig()
	initLogger()
}

func initLogger() {

	cfg, _ := ctx.Value("config").(v1.TortillaConfig)

	logLevelStr := cfg.LogLevel
	level, err := zerolog.ParseLevel(strings.ToLower(logLevelStr))
	if err != nil {
		level = zerolog.InfoLevel
	}
	zerolog.SetGlobalLevel(level)
}

func initConfig() {
	data, err := os.ReadFile(configFileName)
	if err != nil {
		zerolog.Ctx(ctx).Warn().Err(err).Msgf("could not open config file %s, using defaults", configFileName)
		return
	}

	cfg := v1.TortillaConfig{}
	yaml.Unmarshal(data, &cfg)

	ctx = context.WithValue(ctx, "config", cfg)
}

func main() {
	if len(os.Args) < 2 {
		zerolog.Ctx(ctx).Err(fmt.Errorf("no command provided: use 'tortilla --help' for more information")).Send()
		os.Exit(1)
	}

	cmd := os.Args[1]
	args := os.Args[1:]

	switch cmd {
	case "version", "--version", "-v":
		fmt.Println(commandVersion)
	case "help", "--help", "-h":
		fmt.Println(commandLong)
	case "wrap":
	default:
		if e := wrap(ctx, args[0:]); e != nil {
			zerolog.Ctx(ctx).Err(e).Send()
			os.Exit(1)
		}
	}
}
