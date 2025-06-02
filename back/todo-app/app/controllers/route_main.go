// Package controllers provides HTTP handlers and server startup logic.
package controllers

import (
	"net/http"
)

// top ハンドラは、ルート ("/") への HTTP リクエストを処理
// テンプレートを使用して HTML レスポンスを生成し、クライアントに返す
func top(w http.ResponseWriter, r *http.Request) {
	_, err := session(w, r)
	if err != nil {
		// generateHTML 関数を呼び出して、指定されたテンプレートを描画
		generateHTML(w, "Hello", "layout", "public_navbar", "top")
	} else {
		http.Redirect(w, r, "/todos", 302)
	}

}

func index(w http.ResponseWriter, r *http.Request) {
	_, err := session(w, r)
	if err != nil {
		http.Redirect(w, r, "/", 302)
	} else {
		// generateHTML 関数を呼び出して、指定されたテンプレートを描画
		generateHTML(w, nil, "layout", "private_navbar", "index")
	}
}
