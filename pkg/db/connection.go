package db

import (
	_ "github.com/go-sql-driver/mysql"
	"github.com/upper/db/v4/adapter/mysql"
)

func Connect(config *Config, options map[string]string) (Database, error) {
	if options == nil {
		options = map[string]string{}
	}
	if _, ok := options["multiStatements"]; !ok {
		options["multiStatements"] = "1"
	}
	if _, ok := options["charset"]; !ok {
		options["charset"] = "utf8mb4"
	}
	if _, ok := options["collation"]; !ok {
		options["collation"] = "utf8mb4_unicode_ci"
	}
	if _, ok := options["parseTime"]; !ok {
		options["parseTime"] = "1"
	}
	sess, err := mysql.Open(mysql.ConnectionURL{
		Host:     config.Host,
		User:     config.User,
		Password: config.Password,
		Database: config.DatabaseName,
		Options:  options,
	})
	if err != nil {
		return nil, err
	}

	return NewDatabase(sess), nil
}
