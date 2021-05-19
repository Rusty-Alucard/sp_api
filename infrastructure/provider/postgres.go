package config

import (
	"database/sql"
	"fmt"
	"os"
	"sync"
	"time"
)

var singletonDB *sql.DB
var lock sync.Mutex

func connect() (*sql.DB, error) {
	lock.Lock()
	defer lock.Unlock()
	if singletonDB != nil {
		return singletonDB
	}

	singletonDB, err := sql.Open(
		"postgres",
		fmt.Sprintf(
			"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
			os.Getenv("PG_HOST"),
			os.Getenv("PG_PORT"),
			os.Getenv("PG_USER"),
			os.Getenv("PG_PASSWORD"),
			os.Getenv("PG_DATABASE"),
		),
	)

	if err != nil {
		return nil, err
	}

	// setting connection pooling
	singletonDB.SetMaxIdleConns(5)  // コネクションプールの最大数
	singletonDB.SetMaxOpenConns(10) // 接続の最大数
	singletonDB.SetConnMaxIdleTime(10 * time.Second)
	singletonDB.SetConnMaxLifeTime(10 * time.Second) // 接続の再利用可能時間

	return singletonDB, err
}
