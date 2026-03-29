package logger

import (
	"fmt"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"os"
	"time"
)

// Config logger initiation config
type Config struct {
	cfg zap.Config // cfg - internal zap config
	set Settings   // set - passed settings
}

// Settings specific logger configuration set
type Settings struct {
	DisableCaller     bool   // DisableCaller - disable annotating logs with the calling function's file name and line number
	DisableStacktrace bool   // DisableStacktrace - disable automatic stacktrace capturing
	Colour            bool   // Colour - adds colour
	Level             uint8  // Level - set log level
	Path              string // Path - set output log path
	LogName           string // LogName - set output log file name
	FileFormat        string // FileFormat - set output log file extension
	Encoding          string // Encoding - set log format encoding (json/console, default - console)
	CallerKey         string // CallerKey - set specific caller key (for json formatting)
	MessageKey        string // MessageKey - set specific message key (for json formatting)
	TimeKey           string // TimeKey - set specific time key (for json formatting)
	LevelKey          string // LevelKey - set specific level key (for json formatting)
}

// NewConfig create a new logger initiation config
func NewConfig(set Settings) Config {
	return Config{cfg: zap.NewDevelopmentConfig(), set: set}
}

// InitLogger initialise new logger by passed config
func (c *Config) InitLogger() error {

	if c.set.Path != "" && c.set.LogName != "" && c.set.FileFormat != "" {
		if err := c.setOutput(); err != nil {
			return err
		}
	}

	if c.set.Colour {
		c.cfg.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	}
	c.cfg.EncoderConfig.EncodeTime = zapcore.RFC3339TimeEncoder
	c.cfg.DisableCaller = c.set.DisableCaller
	c.cfg.DisableStacktrace = c.set.DisableStacktrace

	c.setLogLevel()

	logger, err := c.cfg.Build()
	if err != nil {
		return err
	}
	zap.ReplaceGlobals(logger)
	logger.Sugar()
	return nil
}

// setOutput set the output settings and create output file for current logger
func (c *Config) setOutput() error {
	output := fmt.Sprintf("%s/%s.%s", c.set.Path, c.set.LogName, c.set.FileFormat)
	if f, err := os.OpenFile(output, os.O_APPEND|os.O_CREATE|os.O_RDWR, 0755); err != nil {
		return err
	} else {
		if _, err = f.WriteString(fmt.Sprintf("Initialised at %v\n", time.Now())); err != nil {
			return err
		}
	}
	c.cfg.OutputPaths = []string{output, "stdout"}
	c.cfg.ErrorOutputPaths = []string{output, "stderr"}
	return nil
}

// setLogLevel set the log level for current logger
func (c *Config) setLogLevel() {
	var level zapcore.Level
	switch c.set.Level {
	case 1:
		level = zap.InfoLevel
	case 2:
		level = zap.WarnLevel
	case 3:
		level = zap.ErrorLevel
	case 4:
		level = zap.FatalLevel
	default:
		level = zap.DebugLevel
	}
	c.cfg.Level.SetLevel(level)
}

// setEncoding set encoding type and add specific names and options
func (c *Config) setEncoding() {
	switch c.set.Encoding {
	case "json":
		c.cfg.Encoding = c.set.Encoding
		c.cfg.EncoderConfig.CallerKey = c.set.CallerKey
		c.cfg.EncoderConfig.MessageKey = c.set.MessageKey
		c.cfg.EncoderConfig.TimeKey = c.set.TimeKey
		c.cfg.EncoderConfig.LevelKey = c.set.LevelKey
	default:
		c.cfg.Encoding = "console"
	}
}
