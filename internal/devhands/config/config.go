package config

import (
	"github.com/metalagman/devhands/pkg/logger"
	"time"
)

type Config struct {
	Server ServerConfig  `mapstructure:"server"`
	Logger logger.Config `mapstructure:"log"`
}

type ServerConfig struct {
	ListenAddr string `mapstructure:"listen"`

	TimeoutRead  time.Duration `mapstructure:"timeout_read"`
	TimeoutWrite time.Duration `mapstructure:"timeout_write"`
	TimeoutIdle  time.Duration `mapstructure:"timeout_idle"`
}
