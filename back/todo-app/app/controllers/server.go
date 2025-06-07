// Package controllers provides HTTP handlers and server startup logic.
package controllers

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"regexp"
	"strconv"
	"todo-app/app/models"
	"todo-app/config"
)

// generateHTML は指定されたテンプレートファイルをパースし、データを適用して HTTP レスポンスライターに書き込む
func generateHTML(w http.ResponseWriter, data interface{}, filenames ...string) {
	var files []string
	for _, file := range filenames {
		// テンプレートファイルのパスを app/views/templates/ ディレクトリからの相対パスとして構築
		files = append(files, fmt.Sprintf("app/views/templates/%s.html", file))
	}

	log.Println("generateHTML: Attempting to parse template files:", files)

	// テンプレートファイルをパース
	// ParseFiles は複数のファイルを読み込み、定義されたテンプレートのセットを作成
	templates, err := template.ParseFiles(files...)
	if err != nil {
		log.Printf("generateHTML: Template parsing error: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	log.Println("generateHTML: Template parsing successful. Executing template with data.")
	// レイアウトテンプレートを基にデータを適用し、レスポンスライターに書き出し
	// ここで "layout" という名前のテンプレートがテンプレートセット内に存在する必要
	err = templates.ExecuteTemplate(w, "layout", data)
	if err != nil {
		log.Printf("generateHTML: Template execution error: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	log.Println("generateHTML: Template executed successfully.")
}

func session(w http.ResponseWriter, r *http.Request) (sess models.Session, err error) {
	cookie, err := r.Cookie("__cookie__")
	if err == nil {
		sess = models.Session{UUID: cookie.Value}
		if ok, _ := sess.CheckSession(); !ok {
			err = fmt.Errorf("Invalid session")
		}
	}
	return sess, err
}

var validPath = regexp.MustCompile("^/todos/(edit|update|delete)/([0-9]+)$")

// parseURL は、URLパスからIDを抽出し、抽出したIDとHTTPレスポンスライター、リクエストオブジェクトを指定されたハンドラ関数に渡す
// パスが正規表現に一致しない場合は 404 Not Found を返す
func parseURL(fn func(http.ResponseWriter, *http.Request, int)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// /todos/edit/1 のようなパスを想定
		q := validPath.FindStringSubmatch(r.URL.Path)
		if q == nil {
			http.NotFound(w, r)
			return
		}
		id, err := strconv.Atoi(q[2])
		if err != nil {
			http.NotFound(w, r)
			return
		}

		fn(w, r, id)
	}
}

// StartMainServer はアプリケーションの Web サーバーを起動し、ルーティングを設定する
func StartMainServer() error {

	// 静的ファイルを提供するためのファイルサーバーを設定
	files := http.FileServer(http.Dir(config.Config.Static))

	// "/static/" パスへのリクエストをファイルサーバーで処理するように設定
	http.Handle("/static/", http.StripPrefix("/static/", files))

	// ルート ("/") へのリクエストを top ハンドラ関数で処理するように設定
	http.HandleFunc("/", top)

	// ルート ("/signup") へのリクエストを signup ハンドラ関数で処理するように設定
	http.HandleFunc("/signup", signup)

	// ルート ("/login") へのリクエストを login ハンドラ関数で処理するように設定
	http.HandleFunc("/login", login)

	// ルート ("/authenticate") へのリクエストを authenticate ハンドラ関数で処理するように設定
	http.HandleFunc("/authenticate", authenticate)

	http.HandleFunc("/logout", logout)

	http.HandleFunc("/todos", index)

	http.HandleFunc("/todos/new", todoNew)

	http.HandleFunc("/todos/save", todoSave)

	// IDを含む /todos/update/{id} 形式のパスを parseURL 経由で todoUpdate ハンドラにルーティング
	http.HandleFunc("/todos/update", parseURL(todoUpdate))
	// IDを含む /todos/edit/{id} 形式のパスを parseURL 経由で todoEdit ハンドラにルーティング
	http.HandleFunc("/todos/edit", parseURL(todoEdit))
	// IDを含む /todos/delete/{id} 形式のパスを parseURL 経由で todoDelete ハンドラにルーティング
	http.HandleFunc("/todos/delete", parseURL(todoDelete))

	// 指定されたポートで HTTP リクエストのリスニングを開始する
	log.Printf("Starting server on port %s...", config.Config.Port) // サーバー起動ログを追加
	return http.ListenAndServe(":"+config.Config.Port, nil)
}
