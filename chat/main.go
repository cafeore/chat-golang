package main

import (
	"flag"
	"github.com/cafeore/chat-golang/trace"
	"log"
	"net/http"
	"os"
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
	//rはtemplateに対して渡す引数(template側で使えるようになる)
	t.templ.Execute(w, r)
}
func main() {
	var addr = flag.String("addr", ":8080", "アプリケーションのアドレス")
	flag.Parse() //フラグを解釈
	r := newRoom()
	r.tracer = trace.New(os.Stdout)
	http.Handle("/chat", MustAuth(&templateHandler{filename: "chat.html"}))
	http.Handle("/login", &templateHandler{filename: "login.html"})
	http.HandleFunc("/auth/", loginHandler)
	http.Handle("/room", r)
	//チャットルームを開始します
	go r.run()
	//Webサーバの開始
	log.Println("Webサーバーを開始します. ポート", *addr)
	if err := http.ListenAndServe(*addr, nil); err != nil {
		log.Fatal("ListenAndServe:", err)
	}
}
