package config

import (
	"context"

	"github.com/vrischmann/envconfig"
)

// Move contains settings for the "move" lambda.
type Move struct {
	Config
}

// LoadMoveConfig reads settings from the environment.
func LoadMoveConfig(ctx context.Context) (*Move, error) {
	cfg := &Move{}
	if err := envconfig.Init(&cfg); err != nil {
		return nil, err
	}
	return cfg, nil
}
