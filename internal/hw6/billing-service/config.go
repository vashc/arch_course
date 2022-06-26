package billing_service

import (
	"arch_course/internal/hw6"
	"github.com/kelseyhightower/envconfig"
)

func NewConfig() (*hw6.Config, error) {
	config := &hw6.Config{}

	if err := envconfig.Process("", config); err != nil {
		return nil, err
	}

	return config, nil
}
