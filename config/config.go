package config

import (
	"flag"
	"fmt"
	"os"
)

var (
	Conf = Config{}
)

type Config struct {
	Database string
	DbUser   string
	DbPass   string
	Dsn      string

	DropboxToken string
}

func init() {
	flag.StringVar(&Conf.Database, "db", "postgres", "database")
	flag.StringVar(&Conf.DbUser, "db-user", "postgres", "database username")
	flag.StringVar(&Conf.DbPass, "db-pass", "postgres", "database password")
	Conf.Dsn = fmt.Sprintf("postgres://%s:%s@postgres:5432/%s?sslmode=disable", Conf.DbUser, Conf.DbPass, Conf.Database)

	Conf.DropboxToken = os.Getenv("DROPBOX_TOKEN")
}
