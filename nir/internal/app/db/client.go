package db

import (
	"context"
	"fmt"
	"os"
	"strconv"

	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/joho/godotenv"
)

type ConfigConnect struct {
	Port     uint64
	Host     string
	DBport   uint64
	User     string
	Password string
	DBname   string
	Sslmode  string
}

func NewDb(ctx context.Context, config ConfigConnect) (*Database, error) {
	pool, err := pgxpool.Connect(ctx, generateDsn(config))
	if err != nil {
		return nil, err
	}
	return newDatabase(pool), nil
}

func generateDsn(config ConfigConnect) string {
	return fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s", config.Host,
		config.Port, config.User, config.Password, config.DBname, config.Sslmode)
}

func LoadEnv(env string) (ConfigConnect, error) {
	err := godotenv.Load(env)
	if err != nil {
		return ConfigConnect{}, err
	}

	dbport, _ := strconv.ParseUint((os.Getenv("dbport")), 10, 64)
	port, _ := strconv.ParseUint((os.Getenv("port")), 10, 64)
	host := os.Getenv("host")
	user := os.Getenv("user")
	password := os.Getenv("password")
	dbname := os.Getenv("dbname")
	sslmode := os.Getenv("sslmode")

	config := ConfigConnect{
		Port:     port,
		Host:     host,
		DBport:   dbport,
		User:     user,
		Password: password,
		DBname:   dbname,
		Sslmode:  sslmode,
	}
	return config, nil
}
