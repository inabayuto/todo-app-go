// Package controllers provides HTTP handlers and server startup logic.
package controllers

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"todo-app/config"
)

// generateHTML は指定されたテンプレートファイルをパースし、データを適用して HTTP レスポンスライターに書き込む
// レイアウトテンプレートと他のテンプレートを組み合わせて使用することを想定しています。
func generateHTML(w http.ResponseWriter, data interface{}, filenames ...string) {
	var files []string
	for _, file := range filenames {
		// テンプレートファイルのパスを app/views/templates/ ディレクトリからの相対パスとして構築
		files = append(files, fmt.Sprintf("app/views/templates/%s.html", file))
	}

	// テンプレートファイルをパース
	// ParseFiles は複数のファイルを読み込み、定義されたテンプレートのセットを作成
	// template.Must は、ParseFiles がエラーを返した場合にパニックを引き起こします。
	templates := template.Must(template.ParseFiles(files...))

	// レイアウトテンプレートを基にデータを適用し、レスポンスライターに書き出し
	// ここで "layout" という名前のテンプレートがテンプレートセット内に存在する必要
	templates.ExecuteTemplate(w, "layout", data)
}

// StartMainServer はアプリケーションの Web サーバーを起動
func StartMainServer() error {

	// 静的ファイルを提供するためのファイルサーバーを設定
	files := http.FileServer(http.Dir(config.Config.Static))

	// "/static/" パスへのリクエストをファイルサーバーで処理するように設定
	http.Handle("/static/", http.StripPrefix("/static/", files))

	// ルート ("/") へのリクエストを top ハンドラ関数で処理するように設定
	http.HandleFunc("/", top)

	// 指定されたポートで HTTP リクエストのリスニングを開始する
	log.Printf("Starting server on port %s...", config.Config.Port) // サーバー起動ログを追加
	return http.ListenAndServe(":"+config.Config.Port, nil)
}
