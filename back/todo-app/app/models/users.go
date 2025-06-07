// モデル定義とDB操作をまとめたパッケージ
package models

import (
	"log"
	"time"
)

// アプリケーションのユーザー情報を保持する構造体
// ID, UUID, 名前, メールアドレス, パスワード, 作成日時, 紐づくTodoリストを持つ
// パスワードはハッシュ化して保存すること
// Todoスライスはユーザーに紐づくタスク一覧
type User struct {
	ID        int       // ユーザーID（主キー）
	UUID      string    // ユーザーごとの一意な識別子
	Name      string    // ユーザー名
	Email     string    // メールアドレス（ログイン用）
	PassWord  string    // ハッシュ化済みパスワード
	CreatedAt time.Time // レコード作成日時
	Todos     []Todo    // ユーザーに紐づくTodoリスト
}

// セッション情報を保持する構造体
// セッションID, UUID, メールアドレス, ユーザーID, 作成日時を持つ
// セッション管理や認証用途で利用
type Session struct {
	ID        int       // セッションID（主キー）
	UUID      string    // セッションごとの一意な識別子
	Email     string    // セッションに紐づくユーザーのメールアドレス
	UserID    int       // 紐づくユーザーID
	CreatedAt time.Time // セッション作成日時
}

// 新規ユーザーをDBに登録する関数
// UUID生成・パスワード暗号化を行い、usersテーブルへINSERT
func (u *User) CreateUser() (err error) {
	cmd := `insert into users (
		uuid,
		name,
		email,
		password,
		created_at) values ($1, $2, $3, $4, $5)`

	_, err = Db.Exec(cmd,
		createUUID(), // UUID生成
		u.Name,
		u.Email,
		Encrypt(u.PassWord), // パスワード暗号化
		time.Now())          // 作成日時

	if err != nil {
		log.Printf("ユーザー作成失敗: %v", err)
		return err
	}
	log.Println("ユーザー作成成功")
	return nil
}

// ユーザーIDでDBからユーザー情報を取得する関数
// 見つからない場合やエラー時はerrを返す
func GetUser(id int) (user User, err error) {
	user = User{}
	cmd := `select id, uuid, name, email, password, created_at
	from users where id = $1`
	err = Db.QueryRow(cmd, id).Scan(
		&user.ID,
		&user.UUID,
		&user.Name,
		&user.Email,
		&user.PassWord,
		&user.CreatedAt,
	)
	return user, err
}

// ユーザー情報（名前・メール）を更新する関数
// IDで該当ユーザーを特定し、name/emailをUPDATE
func (u *User) UpdateUser() (err error) {
	cmd := `update users set name = $1, email = $2 where id = $3`
	_, err = Db.Exec(cmd, u.Name, u.Email, u.ID)
	if err != nil {
		log.Fatalln(err)
	}
	return err
}

// ユーザーIDでDBからユーザーを削除する関数
func (u *User) DeleteUser() (err error) {
	cmd := `delete from users where id = $1`
	_, err = Db.Exec(cmd, u.ID)
	if err != nil {
		log.Fatalln(err)
	}
	return err
}

// メールアドレスでユーザー情報を取得する関数
// 見つからない場合やエラー時はerrを返す
func GetUserByEmail(email string) (user User, err error) {
	user = User{}
	cmd := `select id, uuid, name, email, password, created_at
	from users where email = $1`
	err = Db.QueryRow(cmd, email).Scan(
		&user.ID,
		&user.UUID,
		&user.Name,
		&user.Email,
		&user.PassWord,
		&user.CreatedAt)
	return user, err
}

// ユーザーに紐づく新規セッションをDBに作成する関数
// UUID生成し、sessionsテーブルへINSERT
func (u *User) CreateSession() (session Session, err error) {
	session = Session{}
	cmd1 := `insert into sessions (
		uuid, 
		email, 
		user_id, 
		created_at) values ($1, $2, $3, $4)`

	_, err = Db.Exec(cmd1, createUUID(), u.Email, u.ID, time.Now())
	if err != nil {
		log.Println(err)
	}

	cmd2 := `select id, uuid, email, user_id, created_at
	 from sessions where user_id = $1 and email = $2`

	err = Db.QueryRow(cmd2, u.ID, u.Email).Scan(
		&session.ID,
		&session.UUID,
		&session.Email,
		&session.UserID,
		&session.CreatedAt)

	return session, err
}

// セッションUUIDが有効かDBで検証する関数
// 有効な場合はSession構造体を更新
func (sess *Session) CheckSession() (valid bool, err error) {
	cmd := `select id, uuid, email, user_id, created_at
	 from sessions where uuid = $1`

	err = Db.QueryRow(cmd, sess.UUID).Scan(
		&sess.ID,
		&sess.UUID,
		&sess.Email,
		&sess.UserID,
		&sess.CreatedAt)

	if err != nil {
		valid = false
		return
	}
	if sess.ID != 0 {
		valid = true
	}
	return valid, err
}

// セッションUUIDでDBからセッションを削除する関数
func (sess *Session) DeleteSessionByUUID() (err error) {
	cmd := `delete from sessions where uuid = $1`
	_, err = Db.Exec(cmd, sess.UUID)
	if err != nil {
		log.Fatalln(err)
	}
	return err
}

// セッションに紐づくユーザー情報を取得する関数
// Session.UserIDを使ってusersテーブルから検索
func (sess *Session) GetUserBySession() (user User, err error) {
	user = User{}
	cmd := `select id, uuid, name, email, created_at FROM users
	where id = $1`
	err = Db.QueryRow(cmd, sess.UserID).Scan(
		&user.ID,
		&user.UUID,
		&user.Name,
		&user.Email,
		&user.CreatedAt)
	return user, err
}
