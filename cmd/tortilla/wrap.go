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

	. "github.com/mbovo/tortilla/v1"
	"github.com/rs/zerolog"
	"github.com/spf13/viper"
)

func wrap(ctx context.Context, args []string) (e error) {

	zerolog.Ctx(ctx).Debug().Any("config", viper.GetViper().AllSettings()).Send()

	tortilla := NewTortilla(ctx, viper.GetViper(), args)

	e = tortilla.Prepare()
	if e != nil {
		return
	}

	e = tortilla.Cook()
	if e != nil {
		return
	}

	return tortilla.Wrap()
}
