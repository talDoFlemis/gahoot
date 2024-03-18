package config

import "github.com/taldoflemis/brain.test/internal/adapters/driven/auth"

type localIdpConfig struct {
	ed25519Seed         string `koanf:"auth.ed25519_seed"`
	issuer              string `koanf:"auth.issuer"`
	audience            string `koanf:"auth.audience"`
	accessTimeInMinutes int    `koanf:"auth.access_time_in_minutes"`
	refreshtimeInHours  int    `koanf:"auth.refresh_time_in_hours"`
}

func NewLocalIDPConfig() (*auth.LocalIdpConfig, error) {
	var out localIdpConfig
	err := k.Unmarshal("", &out)
	if err != nil {
		return nil, err
	}
	cfg := auth.NewLocalIdpConfig(
		out.ed25519Seed,
		out.issuer,
		out.audience,
		out.accessTimeInMinutes,
		out.refreshtimeInHours,
	)
	return cfg, nil
}
