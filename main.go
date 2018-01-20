package main

import (
	"log"
	"net/http"
	"path/filepath"
	"sync"
	"text/template"
)

//templは１つのテンプレートを表す
type templateHandler struct {
	once     sync.Once
	filename string
	templ    *template.Template
}

//ServeHTTPはHTTPリクエストを処理する
//既存のstructや型に対して、ServeHTTPメソッドを用意することでhttp.Handleに登録出来るようにする(https://qiita.com/taizo/items/bf1ec35a65ad5f608d45)参照
func (t *templateHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	t.once.Do(func() {
		t.templ = template.Must(template.ParseFiles(filepath.Join("templates", t.filename)))
	})
	t.templ.Execute(w, nil)
}
func main() {
	//ルート
	http.Handle("/", &templateHandler{filename: "chat.html"})
	//Webサーバの開始
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal("ListenAndServe:", err)
	}
}
