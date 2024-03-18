package config

import (
	"github.com/knadh/koanf/parsers/json"
	"github.com/knadh/koanf/providers/file"
	"github.com/knadh/koanf/v2"

	"github.com/taldoflemis/brain.test/internal/ports"
)

var k = koanf.New(".")

type Koanfson struct {
	logger ports.Logger
}

func NewKoanfson(logger ports.Logger) *Koanfson {
	return &Koanfson{
		logger: logger,
	}
}

func (kson *Koanfson) LoadFromJSON(path string) {
	parser := json.Parser()
	if err := k.Load(file.Provider("config.json"), parser); err != nil {
		kson.logger.Error("error loading config: %v", err)
	}
}
