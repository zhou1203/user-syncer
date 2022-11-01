package db

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/mysql"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq"
	"k8s.io/klog/v2"
)

func Migrate(config *Config) error {
	db, err := WaitForConnect(config, 3*time.Minute)
	if err != nil {
		klog.Errorf("waitFor connect failed, %+v, %s", db, err.Error())
		return err
	}

	driver, err := mysql.WithInstance(db, &mysql.Config{})
	if err != nil {
		klog.Errorf("new database driver failed, %s", err.Error())
		return err
	}
	m, err := migrate.NewWithDatabaseInstance(
		"file:///migrate/table",
		"mysql", driver)
	if err != nil {
		klog.Error(fmt.Errorf("new database instance failed, %s", err.Error()))
		return err
	}
	err = m.Up()
	if err != nil {
		if err.Error() != "no change" {
			klog.Error(fmt.Errorf("migrate database failed, %s", err.Error()))
			return err
		}

	}
	klog.Info("migrate database successfully")
	return nil
}

func WaitForConnect(config *Config, timeout time.Duration) (database *sql.DB, err error) {

	klog.Info("start connect database")
	source := fmt.Sprintf("%s:%s@tcp(%s)/%s", config.User, config.Password, config.Host, config.DatabaseName)

	database, err = sql.Open("mysql", source)
	if err == nil && database != nil {
		klog.Info("connect database successfully")
		return
	}

	period := time.Duration(timeout.Seconds()/10) * time.Second
	ticker := time.NewTicker(period)
	defer ticker.Stop()
	done := make(chan bool)

	go func() {
		time.Sleep(timeout)
		done <- true
	}()

	for {
		select {
		case <-ticker.C:
			err = nil
			database, err = sql.Open("mysql", source)
			if err == nil && database != nil {
				klog.Info("connect database successfully")
				return
			} else {
				klog.Info("retry connect database ...")
			}
		case <-done:
			return
		}
	}
}
