package v1

import (
	"context"
	"strings"

	vaultFacade "github.com/mbovo/tortilla/v1/vault"
	"github.com/rs/zerolog"
)

const (
	TransformationTypeRename = "rename"
	TransformationTypePrefix = "prefix"
)

var (
	transformMap = map[string]func(ctx context.Context, tcfg TransformationConfig, template *vaultFacade.GeneratedConfigEnvTemplate) bool{
		TransformationTypeRename: applyRename,
		TransformationTypePrefix: applyPrefix,
	}
)

type Transformer interface {
	Apply(ctx context.Context, vaultConfig *vaultFacade.GeneratedConfig) (*vaultFacade.GeneratedConfig, error)
}

type transformer struct {
	configs []TransformationConfig
}

func applyRename(ctx context.Context, tcfg TransformationConfig, template *vaultFacade.GeneratedConfigEnvTemplate) bool {

	logger := zerolog.Ctx(ctx)
	if template.Name == tcfg.Match {
		logger.Trace().Str("oldName", template.Name).Str("newName", tcfg.Change).Msg("Renaming template")
		template.Name = tcfg.Change
		return true
	}
	return false
}

func applyPrefix(ctx context.Context, tcfg TransformationConfig, template *vaultFacade.GeneratedConfigEnvTemplate) bool {

	logger := zerolog.Ctx(ctx)

	if template.Name == tcfg.Match {
		logger.Trace().Str("prefix", tcfg.Change).Str("templateName", template.Name).Msg("Adding prefix to template")
		template.Name = tcfg.Change + template.Name
		return true
	}

	return false
}

func NewSimpleTransformer(configs []TransformationConfig) Transformer {
	return &transformer{
		configs: configs,
	}
}

func (t *transformer) Apply(ctx context.Context, vaultConfig *vaultFacade.GeneratedConfig) (*vaultFacade.GeneratedConfig, error) {

	logger := zerolog.Ctx(ctx)

	for k, template := range vaultConfig.EnvTemplates {
		logger.Trace().Str("templateName", template.Name).Msg("Checking transformation")

		// for each template check all transformations
		for _, transformation := range t.configs {
			// if path is specified, skip if it doesn't match
			if transformation.Path != "" && !strings.Contains(template.Contents, transformation.Path) {
				logger.Trace().Msg("Skipping transformation due to path mismatch")
				continue
			}

			logger.Trace().
				Str("transformationType", transformation.Type).
				Str("match", transformation.Match).
				Str("change", transformation.Change).
				Str("path", transformation.Path).
				Msg("Evaluating ")

			tFunc, ok := transformMap[transformation.Type]
			if !ok {
				logger.Warn().Str("transformationType", transformation.Type).Msg("Unknown transformation type")
				continue
			}
			changed := tFunc(ctx, transformation, &template)
			if changed {
				logger.Trace().Interface("envTemplate", template).Int("k", k).Msg("Transformed template")
				vaultConfig.EnvTemplates[k] = template
			}

		}
	}
	return vaultConfig, nil
}
