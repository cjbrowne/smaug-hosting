package database

import (
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"github.com/sirupsen/logrus"
	"os"
)

var Connection *sqlx.DB

func Setup() {
	logrus.Trace("Setting up database...")
	conn, err := sqlx.Open("mysql", os.Getenv("MYSQL_CONN_STR"))

	if err != nil {
		logrus.Fatalf("Could not connect to MYSQL Database: %s", err)
	}

	conn.SetMaxIdleConns(2)
	conn.SetMaxOpenConns(95)

	err = conn.Ping()
	if err != nil {
		logrus.Fatalf("MySQL did not respond to PING: %s", err)
	}

	Connection = conn
	logrus.Trace("Database setup")
}