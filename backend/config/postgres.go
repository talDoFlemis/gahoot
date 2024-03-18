package config

import "github.com/taldoflemis/brain.test/internal/adapters/driven/postgres"

type postgresConfig struct {
	Host     string `koanf:"postgres.host"`
	Port     int    `koanf:"postgres.port"`
	User     string `koanf:"postgres.user"`
	Password string `koanf:"postgres.password"`
	Database string `koanf:"postgres.database"`
}

func NewPostgresConfig() (*postgres.Config, error) {
	var out postgresConfig
	err := k.Unmarshal("", &out)
	if err != nil {
		return nil, err
	}
	return postgres.NewConfig(out.User, out.Password, out.Host, out.Database, out.Port), nil
}
