// Package models provides database models and utility functions.
package models

import (
	"crypto/sha1"
	"database/sql"
	"fmt"
	"log"
	"todo-app/config"

	"github.com/google/uuid"
	_ "github.com/lib/pq" // PostgreSQL ドライバーをインポート: database/sql パッケージに PostgreSQL ドライバーを登録するために必要
)

// Db はデータベース接続プールを表すグローバル変数
var Db *sql.DB

// err はデータベース関連のエラーを保持する変数
var err error

// テーブル名の定数
const (
	tableNameUser = "users"
	tableNameTodo = "todos"
	// 他のテーブル名定数もここに追加できます
)

// init 関数はパッケージロード時に自動的に実行されます。
// ここでデータベース接続の初期化とテーブルのセットアップを行います。
func init() {

	// config パッケージがinit関数でLoadConfigを呼び出すように修正済み
	// ここでは既に読み込まれた config.Config を使用して接続文字列を生成します。
	connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		config.Config.DbHost,     // config.Config から取得: データベースホスト名
		config.Config.DbPort,     // config.Config から取得: データベースポート番号
		config.Config.DbUser,     // config.Config から取得: データベースユーザー名
		config.Config.DbPassword, // config.Config から取得: データベースパスワード
		config.Config.DbName)     // config.Config から取得: データベース名

	// データベースに接続します。
	// sql.Open はすぐに接続を確立するわけではなく、DB オブジェクトを準備します。
	Db, err = sql.Open(config.Config.SQLDriver, connStr) // SQLDriverもconfigから取得
	if err != nil {
		log.Fatalln("Failed to open database connection:", err) // 接続失敗時は致命的なエラーとして終了
	}

	// 実際のデータベース接続が確立可能か確認します。
	err = Db.Ping()
	if err != nil {
		log.Fatalln("Failed to connect to database:", err) // 接続確認失敗時も致命的なエラーとして終了
	}
	log.Println("Database connection established successfully!")

	// 開発中はテーブルを自動作成すると都合が良い場合がありますが、
	// 本番環境ではマイグレーションツールなどを使用するのが一般的です。
	// 現在はコメントアウトされていますが、必要に応じてコメントを外してください。
	/*
		// CreateUsersTable は users テーブルが存在しない場合に作成します。
		cmdU := fmt.Sprintf(`CREATE TABLE IF NOT EXISTS %s (
				id SERIAL PRIMARY KEY,
				uuid UUID NOT NULL UNIQUE, // データ型を UUID に修正
				name VARCHAR(255),
				email VARCHAR(255),
				password VARCHAR(255),
				created_at TIMESTAMP
				);`, tableNameUser)

		_, err := Db.Exec(cmdU) // テーブル作成クエリを実行
		if err != nil {
			log.Printf("Error creating %s table: %v", tableNameUser, err) // エラーをログ出力
			// init 関数内でエラーが発生した場合、通常は回復不能なので Fatallnでも良いですが、
			// ここではログ出力に留めます。必要に応じて Fatalln に変更してください。
		}

		log.Printf("%s table creation attempted.", tableNameUser) // 実行試行のログを追加
	*/
	/*
		// CreateTodosTable は todos テーブルが存在しない場合に作成します。
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
}

// createUUID は新しい UUID を生成するヘルパー関数
func createUUID() (uuidobj uuid.UUID) {
	uuidobj, _ = uuid.NewUUID()
	return uuidobj
}

// Encrypt は SHA1 を使用して文字列をハッシュ化する関数 (パスワードなどに使用)
func Encrypt(plaintext string) (cryptext string) {
	cryptext = fmt.Sprintf("%x", sha1.Sum([]byte(plaintext)))
	return cryptext
}
