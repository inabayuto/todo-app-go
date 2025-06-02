package models

import (
	"crypto/sha1"
	"database/sql"
	"fmt"
	"log"
	"todo-app/config"

	"github.com/google/uuid"
	_ "github.com/lib/pq" // PostgreSQL ドライバーをインポート
)

var Db *sql.DB

var err error

const (
	tableNameUser = "users"
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

	// CreateUsersTable は users テーブルが存在しない場合に作成
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

func createUUID() (uuidobj uuid.UUID) {
	uuidobj, _ = uuid.NewUUID()
	return uuidobj
}

func Encrypt(plaintext string) (cryptext string) {
	cryptext = fmt.Sprintf("%x", sha1.Sum([]byte(plaintext)))
	return cryptext
}
