// Package models はデータベースモデルとユーティリティ関数を提供します。
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

// Db はアプリケーション全体で使用されるデータベース接続
var Db *sql.DB

// err はデータベース関連のエラーを保持する変数
var err error

// テーブル名の定数
const (
	tableNameUser    = "users"
	tableNameTodo    = "todos"
	tableNameSession = "sessions"
)

// ここでデータベース接続の初期化とテーブルのセットアップを行います。
func init() {

	// データベース接続文字列を生成します。
	connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		config.Config.DbHost,     // データベースホスト名
		config.Config.DbPort,     // データベースポート番号
		config.Config.DbUser,     // データベースユーザー名
		config.Config.DbPassword, // データベースパスワード
		config.Config.DbName)     // データベース名

	// データベースに接続を試みます。
	// sql.Open はすぐに接続を確立するのではなく、DB オブジェクトを準備します。
	Db, err = sql.Open(config.Config.SQLDriver, connStr) // SQLDriverもconfigから取得
	if err != nil {
		log.Fatalln("Failed to open database connection:", err) // 接続失敗時は致命的なエラーとして終了
	}

	// データベースへの接続を確認します。
	err = Db.Ping()
	if err != nil {
		log.Fatalln("Failed to connect to database:", err) // 接続確認失敗時も致命的なエラーとして終了
	}
	log.Println("Database connection established successfully!") // 接続成功をログ出力

	/*
		// ユーザーテーブルを作成するSQLコマンド
		cmdU := fmt.Sprintf(`CREATE TABLE IF NOT EXISTS %s (
				id SERIAL PRIMARY KEY,
				uuid UUID NOT NULL UNIQUE,
				name VARCHAR(255),
				email VARCHAR(255),
				password VARCHAR(255),
				created_at TIMESTAMP
				);`, tableNameUser)

		_, err := Db.Exec(cmdU) // テーブル作成クエリを実行
		if err != nil {
			log.Printf("Error creating %s table: %v", tableNameUser, err) // エラーをログ出力
		}

		log.Printf("%s table creation attempted.", tableNameUser) // 実行試行のログを追加
	*/
	/*
		// Todoテーブルを作成するSQLコマンド
		cmdT := fmt.Sprintf(`CREATE TABLE IF NOT EXISTS %s(
				id SERIAL PRIMARY KEY, // AUTOINCREMENT から SERIAL に変更
				content TEXT,
				user_id INTEGER,
				created_at TIMESTAMP
				);`, tableNameTodo)

		_, err = Db.Exec(cmdT) // テーブル作成クエリを実行
		if err != nil {
			log.Printf("Error creating %s table: %v", tableNameTodo, err) // エラーをログ出力
		}

		log.Printf("%s table creation attempted.", tableNameTodo) // 実行試行のログを追加
	*/

	// セッションテーブルを作成するSQLコマンド
	cmdS := fmt.Sprintf(`CREATE TABLE IF NOT EXISTS %s(
			id SERIAL PRIMARY KEY,
			uuid UUID NOT NULL UNIQUE,
			email VARCHAR(255),
			user_id INTEGER,
			created_at TIMESTAMP)`, tableNameSession)

	_, err := Db.Exec(cmdS) // テーブル作成クエリを実行
	if err != nil {
		log.Printf("Error creating %s table: %v", tableNameSession, err) // エラーをログ出力
	}

	log.Printf("%s table creation attempted.", tableNameSession) // 実行試行のログを追加
}

// createUUID は新しいUUIDを生成するヘルパー関数
func createUUID() (uuidobj uuid.UUID) {
	uuidobj, _ = uuid.NewUUID()
	return uuidobj
}

// Encrypt は SHA1 を使用して文字列をハッシュ化する関数
func Encrypt(plaintext string) (cryptext string) {
	cryptext = fmt.Sprintf("%x", sha1.Sum([]byte(plaintext)))
	return cryptext
}
