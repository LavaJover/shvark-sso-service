package logger

import (
	"github.com/LavaJover/shvark-sso-service/internal/config"
	"go.uber.org/zap"
)

var log *zap.Logger

func InitLogger(cfg *config.SSOConfig) (*zap.Logger, error) {
	
}