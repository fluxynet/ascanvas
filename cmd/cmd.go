package cmd

import (
	"encoding/json"
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// DoNotLogToFile is an indication not to log to any file
const DoNotLogToFile = ""

// Logger is a factory for a logger.
func Logger(level string, filename string) (*zap.Logger, error) {
	var cfg = zap.NewProductionConfig()
	cfg.OutputPaths = []string{
		"stdout",
	}

	if filename != "" {
		_ = os.Truncate(filename, 0)
		cfg.OutputPaths = append(cfg.OutputPaths, filename)
	}

	var lvl zapcore.Level
	if err := lvl.Set(level); err != nil {
		return nil, err
	}

	cfg.Level = zap.NewAtomicLevelAt(lvl)
	cfg.EncoderConfig.EncodeTime = zapcore.RFC3339TimeEncoder
	cfg.DisableCaller = true
	// cfg.DisableStacktrace = true

	return cfg.Build()
}

// LoadConfig from json file into target
func LoadConfig(filename string, target interface{}) error {
	b, err := os.ReadFile(filename)

	if err != nil {
		return err
	}

	err = json.Unmarshal(b, target)

	return err
}
