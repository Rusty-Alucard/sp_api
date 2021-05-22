package provider

import (
	"database/sql"
	"fmt"
	"sync"
	"time"

	_ "github.com/lib/pq"

	"github.com/Rusty-Alucard/sp_api/config"
)

var (
	instance *sql.DB
	lock     sync.Mutex
)

func Connect() (*sql.DB, error) {
	lock.Lock()
	defer lock.Unlock()

	if instance != nil {
		return instance, nil
	}

	db, err := sql.Open(
		config.GetConfig().Database.Driver,
		fmt.Sprintf(
			"host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
			config.GetConfig().Database.Host,
			config.GetConfig().Database.Port,
			config.GetConfig().Database.User,
			config.GetConfig().Database.Password,
			config.GetConfig().Database.Db,
		),
	)
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

	// setting connection pooling
	db.SetMaxIdleConns(5)  // コネクションプールの最大数
	db.SetMaxOpenConns(10) // 接続の最大数
	db.SetConnMaxIdleTime(10 * time.Second)
	db.SetConnMaxLifetime(10 * time.Second) // 接続の再利用可能時間

	instance = db

	return instance, nil
}
