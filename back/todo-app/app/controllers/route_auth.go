package controllers

import (
	"log"
	"net/http"
	"todo-app/app/models"
)

// ハンドラ
// w：http.ResponseWriter クライアントへのレスポンスを書き込むためのオブジェクト（出力）
// r：*http.Request クライアントから送られてきたリクエスト情報（入力）

// signupハンドラ: ユーザーサインアップ処理を担当
// GET: サインアップフォーム表示、POST: ユーザー登録処理
func signup(w http.ResponseWriter, r *http.Request) {
	// リクエストメソッドで分岐
	if r.Method == "GET" {
		_, err := session(w, r)
		if err != nil {
			// 未ログイン時はサインアップフォームを表示
			generateHTML(w, nil, "layout", "signup", "public_navbar")
		} else {
			// ログイン済みならTodo一覧へリダイレクト
			http.Redirect(w, r, "/todos", http.StatusFound)
		}
	} else if r.Method == "POST" {
		// フォーム送信時の処理
		err := r.ParseForm()
		if err != nil {
			log.Println("Form parsing error:", err)
		}
		// フォーム値からユーザー情報を生成
		user := models.User{
			Name:     r.PostFormValue("name"),
			Email:    r.PostFormValue("email"),
			PassWord: r.PostFormValue("password"),
		}
		// DBにユーザー登録
		if err := user.CreateUser(); err != nil {
			log.Println("Database user creation error:", err)
		} else {
			log.Println("User created successfully.")
			http.Redirect(w, r, "/", http.StatusFound)
		}
	}
}

// loginハンドラ: ログインフォーム表示のみ担当
func login(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		_, err := session(w, r)
		if err != nil {
			generateHTML(w, nil, "layout", "login", "public_navbar")
		} else {
			http.Redirect(w, r, "/todos", http.StatusFound)
		}
	}
}

// authenticateハンドラ: ログイン認証処理を担当
func authenticate(w http.ResponseWriter, r *http.Request) {
	log.Println("authenticate handler started")
	// フォームデータをパース。POSTで送信された値を扱うため必須
	err := r.ParseForm()
	if err != nil {
		log.Println("Form parse error in authenticate:", err)
		// フォームパース失敗時は400エラーを返して終了
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	// 入力されたメールアドレスでユーザーをDBから検索
	log.Println("Attempting to get user by email:", r.PostFormValue("email"))
	user, err := models.GetUserByEmail(r.PostFormValue("email"))
	if err != nil {
		// ユーザーが見つからない場合はログイン画面へリダイレクト
		log.Println("Error getting user by email:", err)
		http.Redirect(w, r, "/login", http.StatusFound)
		return
	}

	// パスワード照合（DB保存値はハッシュ化済みなので同じ関数で暗号化して比較）
	log.Println("User found, comparing passwords.")
	if user.PassWord == models.Encrypt(r.PostFormValue("password")) {
		// パスワード一致時はセッション作成
		log.Println("Password matched. Creating session.")
		session, err := user.CreateSession()
		if err != nil {
			// セッション作成失敗時はログイン画面へリダイレクト
			log.Println("Error creating session:", err)
			http.Redirect(w, r, "/login", http.StatusFound)
			return
		}

		// セッションUUIDをクッキーに保存（HttpOnlyでJSからアクセス不可）
		log.Println("Session created. Setting cookie.")
		cookie := http.Cookie{
			Name:     "__cookie__",
			Value:    session.UUID,
			HttpOnly: true,
		}
		http.SetCookie(w, &cookie)

		// 認証成功後はTodo一覧へリダイレクト
		log.Println("Cookie set. Redirecting to /todos.")
		http.Redirect(w, r, "/todos", http.StatusFound)
	} else {
		// パスワード不一致時はログイン画面へリダイレクト
		log.Println("Incorrect password for email:", r.PostFormValue("email"))
		http.Redirect(w, r, "/login", http.StatusFound)
	}
}

// logoutハンドラ: ログアウト処理を担当
func logout(w http.ResponseWriter, r *http.Request) {
	// クッキーからセッションUUIDを取得。未ログイン時はエラーになる
	cookie, err := r.Cookie("__cookie__")
	if err != nil {
		log.Println(err)
	}

	if err != http.ErrNoCookie {
		// セッションUUIDが存在する場合はDBから該当セッションを削除
		session := models.Session{UUID: cookie.Value}
		session.DeleteSessionByUUID()
	}

	// セッションクッキーを無効化（MaxAge=-1で即時削除）
	http.SetCookie(w, &http.Cookie{
		Name:     "__cookie__",
		Value:    "",
		Path:     "/",
		MaxAge:   -1,
		HttpOnly: true,
	})

	// ログアウト後はログイン画面へリダイレクト
	http.Redirect(w, r, "/login", http.StatusFound)
}
