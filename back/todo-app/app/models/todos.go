package models

import (
	"log"
	"time"
)

// Todo構造体はアプリケーションの単一のTodoアイテムを表す
// TodoのID、内容、所有ユーザーのID、作成日時を含む
type Todo struct {
	ID        int       // Todoアイテムの一意なID
	Content   string    // Todoの内容
	UserID    int       // このTodoを所有するユーザーのID
	CreatedAt time.Time // Todoが作成された日時
}

// Todoの内容を入力として受け取り、呼び出し元のUser構造体のIDに関連づける
func (u *User) CreateTodo(content string) (err error) {
	// 新しいTodoをtodosテーブルに挿入するSQLコマンド
	cmd := `insert into todos (
		content,
		user_id,
		created_at) values ($1, $2, $3)`

	// SQLコマンドを実行し、Todo内容、ユーザーID、現在時刻を挿入
	_, err = Db.Exec(cmd, content, u.ID, time.Now())
	if err != nil {
		// 実行失敗した場合に致命的なエラーをログ出力
		log.Fatalln(err)
	}
	return err
}

// IDを指定してデータベースから単一のTodoアイテムを取得
// Todo構造体と、取得に失敗した場合のエラーを返す
func GetTodo(id int) (todo Todo, err error) {
	// IDを指定してtodosテーブルからTodoを取得するSQLコマンド
	cmd := `select id, content, user_id, created_at from todos
	where id = $1`
	// 空のTodo構造体を初期化
	todo = Todo{}

	// クエリを実行し、結果をtodo構造体のフィールドにスキャン
	err = Db.QueryRow(cmd, id).Scan(
		&todo.ID,
		&todo.Content,
		&todo.UserID,
		&todo.CreatedAt)

	// 取得したTodoとエラーを返す
	return todo, err
}

// データベースから全てのTodoアイテムを取得します
// Todo構造体のスライスと、取得に失敗した場合のエラーを返します
func GetTodos() (todos []Todo, err error) {
	// 全てのTodoをtodosテーブルから取得するSQLコマンド
	cmd := `select id, content, user_id, created_at from todos`
	// クエリを実行して全ての行を取得
	rows, err := Db.Query(cmd)
	if err != nil {
		// クエリ失敗した場合に致命的なエラーをログ出力
		log.Fatalln(err)
	}
	// リソース解放のため、関数の最後にrows.Close()が実行されるように遅延設定
	defer rows.Close()

	// 行をイテレート
	for rows.Next() {
		var todo Todo // 各行のTodo構造体を宣言
		// 行データをTodo構造体のフィールドにスキャン
		err = rows.Scan(&todo.ID,
			&todo.Content,
			&todo.UserID,
			&todo.CreatedAt)
		if err != nil {
			// スキャン失敗した場合に致命的なエラーをログ出力
			log.Fatalln(err)
		}
		// スキャンしたtodoをスライスに追加
		todos = append(todos, todo)
	}
	// rows.Close() // deferで実行されるためコメントアウト

	// Todoのスライスとエラーを返す
	return todos, err
}

// 特定のユーザーに紐づく全てのTodoアイテムを取得
// User構造体を引数として受け取り、Todo構造体のスライスとエラーを返す
func (u *User) GetTodosByUser() (todos []Todo, err error) {
	// 特定のユーザーIDでtodosテーブルからTodoを取得するSQLコマンド
	cmd := `select id, content, user_id, created_at from todos
	where user_id = $1`

	// 特定のユーザーIDでクエリを実行
	rows, err := Db.Query(cmd, u.ID)
	if err != nil {
		// クエリ失敗した場合に致命的なエラーをログ出力
		log.Fatalln(err)
	}
	// リソース解放のため、関数の最後にrows.Close()が実行されるように遅延設定
	defer rows.Close()

	// 行をイテレートします。
	for rows.Next() {
		var todo Todo // 各行のTodo構造体を宣言
		// 行データをTodo構造体のフィールドにスキャン
		err = rows.Scan(
			&todo.ID,
			&todo.Content,
			&todo.UserID,
			&todo.CreatedAt)

		if err != nil {
			// スキャン失敗した場合に致命的なエラーをログ出力
			log.Fatalln(err)
		}
		// スキャンしたtodoをスライスに追加
		todos = append(todos, todo)
	}
	// rows.Close() // deferで実行されるためコメントアウト

	// Todoのスライスとエラーを返す
	return todos, err
}

// データベース内の既存のTodoアイテムを更新
// 呼び出し元のTodo構造体のIDを使用して、更新するアイテムを特定
func (t *Todo) UpdateTodo() error {
	// Todo情報を更新するSQLコマンド
	cmd := `update todos set content = $1, user_id = $2
	where id = $3`
	// 新しい内容、ユーザーID、Todo IDで更新コマンドを実行
	_, err = Db.Exec(cmd, t.Content, t.UserID, t.ID)
	if err != nil {
		// エラーをログ出力
		log.Printf("Error updating todo (ID %d): %v", t.ID, err)
		return err
	}
	// 成功をログ出力
	log.Printf("Successfully updated todo (ID %d)", t.ID)
	return nil
}

// IDを指定してデータベースからTodoアイテムを削除
// 呼び出し元のTodo構造体のIDを使用して、削除するアイテムを特定
func (t *Todo) DeleteTodo() error {
	// Todoを削除するSQLコマンド
	cmd := `delete from todos where id = $1`
	// Todo IDで削除コマンドを実行
	_, err = Db.Exec(cmd, t.ID)
	if err != nil {
		// エラーをログ出力
		log.Printf("Error deleting todo (ID %d): %v", t.ID, err)
		return err
	}
	// 成功をログ出力
	log.Printf("Successfully deleted todo (ID %d)", t.ID)
	return nil
}
