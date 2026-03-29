package logger

import "testing"

func TestInitLogger(t *testing.T) {
	set := Settings{
		DisableCaller:     false,
		DisableStacktrace: false,
		Colour:            false,
		Level:             0,
		Path:              "",
		LogName:           "",
		FileFormat:        "",
		Encoding:          "",
		CallerKey:         "",
		MessageKey:        "",
		TimeKey:           "",
		LevelKey:          "",
	}
	cfg := NewConfig(set)
	if err := cfg.InitLogger(); err != nil {
		t.Errorf("failed to init logger %v", err)
	}
}
