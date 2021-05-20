package config

import (
	"database/sql"
	"fmt"
	"sync"
	"time"

	"github.com/Rusty-Alucard/sp_api/config"
)

var singletonDB *sql.DB
var lock sync.Mutex

func connect(cfg Config) (*sql.DB, error) {
	lock.Lock()
	defer lock.Unlock()
	if singletonDB != nil {
		return singletonDB, nil
	}

	singletonDB, err := sql.Open(
		"postgres",
		fmt.Sprintf(
			"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
			cfg.Host,
			cfg.Port,
			cfg.User,
			cfg.Password,
			cfg.Db,
		),
	)

	if err != nil {
		return nil, err
	}

	// setting connection pooling
	singletonDB.SetMaxIdleConns(5)  // コネクションプールの最大数
	singletonDB.SetMaxOpenConns(10) // 接続の最大数
	singletonDB.SetConnMaxIdleTime(10 * time.Second)
	singletonDB.SetConnMaxLifetime(10 * time.Second) // 接続の再利用可能時間

	return singletonDB, err
}
