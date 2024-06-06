package config

import (
	"encoding/json"
	"errors"
	"fmt"
)

func Parse(configType Type, configData []byte) (*Config, error) {
	var cfg Config

	switch configType {
	case JSON:
		err := json.Unmarshal(configData, &cfg)
		if err != nil {
			return nil, fmt.Errorf("failed to parse json config: %w", err)
		}
	case Swagger:
		return nil, errors.New("swagger config is not supported yet")
	default:
		return nil, fmt.Errorf("unknown config type: %s", configType)
	}

	return &cfg, nil
}
