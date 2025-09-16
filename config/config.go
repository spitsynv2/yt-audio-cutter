package config

import (
	"flag"
	"fmt"
)

var (
	Conf = Config{}
)

type Config struct {
	Database string
	DbUser   string
	DbPass   string

	Dsn string
}

func init() {
	flag.StringVar(&Conf.Database, "db", "postgres", "database")
	flag.StringVar(&Conf.DbUser, "db-user", "yt", "database username")
	flag.StringVar(&Conf.DbPass, "db-pass", "postgres", "database password")

	Conf.Dsn = fmt.Sprintf("postgres://%s:%s@localhost:5432/%s?sslmode=disable", Conf.DbUser, Conf.DbPass, Conf.Database)
}
