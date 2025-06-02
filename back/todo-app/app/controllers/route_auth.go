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
		_, err := session(w, r)
		if err != nil {
			// GET リクエストの場合はサインアップフォームのページを表示
			generateHTML(w, nil, "layout", "signup", "public_navbar")
		} else {
			http.Redirect(w, r, "/todos", 302)
		}
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
		} else {
			// ユーザー作成成功時の処理
			log.Println("User created successfully.") // 成功ログを追加
			http.Redirect(w, r, "/", http.StatusFound)
		}
	}
}

func login(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		_, err := session(w, r)
		if err != nil {
			generateHTML(w, nil, "layout", "login", "public_navbar")
		} else {
			http.Redirect(w, r, "/todos", 302)
		}
	}
}

func authenticate(w http.ResponseWriter, r *http.Request) {
	log.Println("authenticate handler started")
	err := r.ParseForm()
	if err != nil {
		log.Println("Form parse error in authenticate:", err)
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	log.Println("Attempting to get user by email:", r.PostFormValue("email"))
	user, err := models.GetUserByEmail(r.PostFormValue("email"))
	if err != nil {
		log.Println("Error getting user by email:", err)
		http.Redirect(w, r, "/login", http.StatusFound)
		return
	}

	log.Println("User found, comparing passwords.")
	if user.PassWord == models.Encrypt(r.PostFormValue("password")) {
		log.Println("Password matched. Creating session.")
		session, err := user.CreateSession()
		if err != nil {
			log.Println("Error creating session:", err)
			http.Redirect(w, r, "/login", http.StatusFound)
			return
		}

		log.Println("Session created. Setting cookie.")
		cookie := http.Cookie{
			Name:     "__cookie__",
			Value:    session.UUID,
			HttpOnly: true,
		}
		http.SetCookie(w, &cookie)

		log.Println("Cookie set. Redirecting to /todos.")
		http.Redirect(w, r, "/todos", http.StatusFound)
	} else {
		log.Println("Incorrect password for email:", r.PostFormValue("email"))
		http.Redirect(w, r, "/login", http.StatusFound)
	}
}

func logout(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("__cookie__")
	if err != nil {
		log.Println(err)
	}

	if err != http.ErrNoCookie {
		session := models.Session{UUID: cookie.Value}
		session.DeleteSessionByUUID()
	}

	// ブラウザからセッションクッキーを削除
	http.SetCookie(w, &http.Cookie{
		Name:     "__cookie__",
		Value:    "",
		Path:     "/", // クッキーのパスを適切に設定
		MaxAge:   -1,  // クッキーを即時削除
		HttpOnly: true,
	})

	http.Redirect(w, r, "/login", http.StatusFound)

}
