package provider

import (
	"database/sql"
	"fmt"
	"sync"
	"time"

	_ "github.com/lib/pq"

	"github.com/Rusty-Alucard/sp_api/config"
)

var singletonDB *sql.DB
var lock sync.Mutex

func Connect() (*sql.DB, error) {
	lock.Lock()
	defer lock.Unlock()

	if singletonDB != nil {
		fmt.Println("--singleton")
		return singletonDB, nil
	}

	singletonDB, err := sql.Open(
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

	if err := singletonDB.Ping(); err != nil {
		return nil, err
	}

	// setting connection pooling
	singletonDB.SetMaxIdleConns(5)  // コネクションプールの最大数
	singletonDB.SetMaxOpenConns(10) // 接続の最大数
	singletonDB.SetConnMaxIdleTime(10 * time.Second)
	singletonDB.SetConnMaxLifetime(10 * time.Second) // 接続の再利用可能時間

	fmt.Println("--init")

	return singletonDB, nil
}
