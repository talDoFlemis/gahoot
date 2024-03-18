package postgres

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	pgxUUID "github.com/vgarvardt/pgx-google-uuid/v5"
)

var (
	globalPool *pgxpool.Pool
)

type Config struct {
	User     string
	Password string
	Host     string
	Port     int
	Database string
}

func NewConfig(user, password, domain, database string, port int) *Config {
	return &Config{
		User:     user,
		Password: password,
		Host:     domain,
		Port:     port,
	}
}

func NewConnection(connStr string) (*pgx.Conn, error) {
	connConfig, err := pgx.ParseConfig(connStr)
	if err != nil {
		return nil, err
	}

	conn, err := pgx.ConnectConfig(context.Background(), connConfig)
	if err != nil {
		return nil, err
	}
	pgxUUID.Register(conn.TypeMap())
	return conn, nil
}

func NewPool(connStr string) (*pgxpool.Pool, error) {
	pgxConfig, err := pgxpool.ParseConfig(connStr)
	if err != nil {
		return nil, err
	}

	pgxConfig.AfterConnect = func(ctx context.Context, conn *pgx.Conn) error {
		pgxUUID.Register(conn.TypeMap())
		return nil
	}

	pool, err := pgxpool.NewWithConfig(context.TODO(), pgxConfig)
	if err != nil {
		return nil, err
	}

	return pool, nil
}

func GenerateConnectionString(cfg *Config) string {
	str := fmt.Sprintf(
		"postgres://%s:%s@%s:%d/%s&sslmode=disable",
		cfg.User,
		cfg.Password,
		cfg.Host,
		cfg.Port,
		cfg.Database,
	)
	return str
}
