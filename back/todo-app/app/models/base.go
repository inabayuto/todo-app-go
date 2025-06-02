package models

import (
	"database/sql"
	"fmt"
	"log"
	"todo-app/config"

	_ "github.com/lib/pq" // PostgreSQL ドライバーをインポート
)

var Db *sql.DB

var err error

const (
	tableNameUser = "users"
	// 他のテーブル名定数もここに追加できます
)

func init() {

	connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		config.Config.DbHost,     // config.Config から取得
		config.Config.DbPort,     // config.Config から取得
		config.Config.DbUser,     // config.Config から取得
		config.Config.DbPassword, // config.Config から取得
		config.Config.DbName)     // config.Config から取得

	Db, err = sql.Open(config.Config.SQLDriver, connStr) // SQLDriverもconfigから取得
	if err != nil {
		log.Fatalln(err)
	}

	// CreateUsersTable は users テーブルが存在しない場合に作成します。
	cmdU := fmt.Sprintf(`CREATE TABLE IF NOT EXISTS %s (
		id SERIAL PRIMARY KEY,
		uuid VARCHAR(255) NOT NULL UNIQUE,
		name VARCHAR(255),
		email VARCHAR(255),
		password VARCHAR(255),
		created_at TIMESTAMP
		);`, tableNameUser)

	_, err := Db.Exec(cmdU)
	if err != nil {
		log.Printf("Error creating %s table: %v", tableNameUser, err)
	}

	log.Printf("%s table creation attempted.", tableNameUser)
}
