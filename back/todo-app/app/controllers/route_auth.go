// Package controllers provides HTTP handlers.
package controllers

import (
	"log"
	"net/http"
	"todo-app/app/models"
)

// signup ハンドラは、ユーザーサインアップのリクエストを処理を実施する
func signup(w http.ResponseWriter, r *http.Request) {
	// リクエストメソッドに応じて処理を分岐
	if r.Method == "GET" {
		// GET リクエストの場合はサインアップフォームのページを表示
		generateHTML(w, nil, "layout", "signup", "public_navbar")
	} else if r.Method == "POST" {
		// POST リクエストの場合はフォームデータを処理

		// フォームデータをパース
		err := r.ParseForm()
		if err != nil {
			log.Println("Form parsing error:", err)

		}

		// フォームからユーザー情報を取得し、新しい User オブジェクトを作成
		user := models.User{
			Name:     r.PostFormValue("name"),     // フォームの "name" フィールドから取得
			Email:    r.PostFormValue("email"),    // フォームの "email" フィールドから取得
			PassWord: r.PostFormValue("password"), // フォームの "password" フィールドから取得
		}

		// ユーザー情報をデータベースに保存
		if err := user.CreateUser(); err != nil {
			log.Println("Database user creation error:", err)

			// データベース保存の成功・失敗に関わらず、ルートパスにリダイレクト
			http.Redirect(w, r, "/", http.StatusFound) // 302 は http.StatusFound
		}
	}
}
