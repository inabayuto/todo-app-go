// Package controllers provides HTTP handlers and server startup logic.
package controllers

import (
	"net/http"
)

// top ハンドラは、ルート ("/") への HTTP リクエストを処理
// テンプレートを使用して HTML レスポンスを生成し、クライアントに返す
func top(w http.ResponseWriter, r *http.Request) {
	// generateHTML 関数を呼び出して、指定されたテンプレートを描画
	generateHTML(w, "Hello", "layout", "public_navbar", "top")
}
