// Package controllers provides HTTP handlers and server startup logic.
package controllers

import (
	"log"
	"net/http"
	"todo-app/app/models"
)

// top ハンドラは、ルート ("/") への HTTP リクエストを処理
// テンプレートを使用して HTML レスポンスを生成し、クライアントに返す
func top(w http.ResponseWriter, r *http.Request) {
	_, err := session(w, r)
	if err != nil {
		// generateHTML 関数を呼び出して、指定されたテンプレートを描画
		generateHTML(w, nil, "layout", "public_navbar", "top")
	} else {
		http.Redirect(w, r, "/todos", http.StatusFound)
	}
}

// index ハンドラは、ユーザーのTodoリストを表示する
// セッションを確認し、認証済みであればユーザーのTodoを取得してテンプレートに渡す
func index(w http.ResponseWriter, r *http.Request) {
	log.Println("index handler started")
	sess, err := session(w, r)
	if err != nil {
		log.Println("index handler: session error, redirecting to /:", err)
		http.Redirect(w, r, "/", http.StatusFound)
	} else {
		log.Println("index handler: session found")
		user, err := sess.GetUserBySession()
		if err != nil {
			log.Println("index handler: Error getting user by session:", err)
		}
		todos, _ := user.GetTodosByUser()
		user.Todos = todos
		log.Printf("index handler: User object before passing to template: %+v\n", user)
		// generateHTML 関数を呼び出して、指定されたテンプレートを描画
		generateHTML(w, user, "layout", "private_navbar", "index")
	}
}

// todoNew ハンドラは、新しいTodo作成フォームを表示する
// 認証済みであればフォームページを表示する
func todoNew(w http.ResponseWriter, r *http.Request) {
	_, err := session(w, r)
	if err != nil {
		http.Redirect(w, r, "/login", http.StatusFound)
	} else {
		// generateHTML 関数を呼び出して、指定されたテンプレートを描画
		generateHTML(w, nil, "layout", "private_navbar", "todo_new")
	}
}

// todoSave ハンドラは、新しいTodoの作成リクエストを処理する
// フォームから内容を取得し、ユーザーに関連付けて保存後、一覧ページにリダイレクトする
func todoSave(w http.ResponseWriter, r *http.Request) {
	log.Println("todoSave handler started")
	sess, err := session(w, r)
	if err != nil {
		log.Println("todoSave handler: session error, redirecting to /login:", err)
		http.Redirect(w, r, "/login", http.StatusFound)
		return
	}

	log.Println("todoSave handler: session found, parsing form.")
	err = r.ParseForm()
	if err != nil {
		log.Println("todoSave handler: Form parse error:", err)
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	log.Println("todoSave handler: Attempting to get user by session.")
	user, err := sess.GetUserBySession()
	if err != nil {
		log.Println("todoSave handler: Error getting user by session:", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	content := r.PostFormValue("content")
	if content == "" {
		log.Println("todoSave handler: Content is empty.")
		http.Error(w, "Bad Request: Todo content cannot be empty", http.StatusBadRequest)
		return
	}

	log.Printf("todoSave handler: Creating todo for user %d with content: %s", user.ID, content)
	if err := user.CreateTodo(content); err != nil {
		log.Println("todoSave handler: Error creating todo:", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	log.Println("todoSave handler: Todo created successfully, redirecting to /todos.")
	http.Redirect(w, r, "/todos", http.StatusFound)
}

// todoEdit ハンドラは、既存のTodoの編集フォームを表示する
// URLパスからTodo IDを取得し、Todo情報を取得してテンプレートに渡す
func todoEdit(w http.ResponseWriter, r *http.Request, id int) {
	sess, err := session(w, r)
	if err != nil {
		http.Redirect(w, r, "/login", 302)
	} else {
		_, err := sess.GetUserBySession()
		if err != nil {
			log.Println(err)
		}
		t, err := models.GetTodo(id)
		if err != nil {
			log.Println(err)
		}
		generateHTML(w, t, "layout", "private_navbar", "todo_edit")
	}
}

// todoUpdate ハンドラは、既存のTodoの更新リクエストを処理する
// URLパスからTodo ID、フォームから更新内容を取得し、Todoを更新後、一覧ページにリダイレクトする
func todoUpdate(w http.ResponseWriter, r *http.Request, id int) {
	sess, err := session(w, r)
	if err != nil {
		http.Redirect(w, r, "/login", 302)
	} else {
		err := r.ParseForm()
		if err != nil {
			log.Println(err)
		}
		user, err := sess.GetUserBySession()
		if err != nil {
			log.Println(err)
		}
		content := r.PostFormValue("content")
		t := &models.Todo{ID: id, Content: content, UserID: user.ID}
		if err := t.UpdateTodo(); err != nil {
			log.Println(err)
		}
		http.Redirect(w, r, "/todos", 302)
	}
}

// todoDelete ハンドラは、既存のTodoの削除リクエストを処理する
// URLパスからTodo IDを取得し、Todoを削除後、一覧ページにリダイレクトする
func todoDelete(w http.ResponseWriter, r *http.Request, id int) {
	sess, err := session(w, r)
	if err != nil {
		http.Redirect(w, r, "/login", http.StatusFound)
	} else {
		_, err := sess.GetUserBySession()
		if err != nil {
			log.Println(err)
		}
		t, err := models.GetTodo(id)
		if err != nil {
			log.Println(err)
		}
		if err := t.DeleteTodo(); err != nil {
			log.Println(err)
		}
		http.Redirect(w, r, "/todos", http.StatusFound)
	}
}
