package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"
	"path"
	"strings"

	"./parser"
	"./websocket_server"
	"./worker"
)

var (
	port     = flag.Int("p", 8080, "Port to use")
	markpath = flag.String("i", "mark.md", "Path of markdown file")
	temppath = flag.String("t", "static/index.html", "Path of template html file.")
)

const ContentsKeyword = "<% contents %>"

func handler(w http.ResponseWriter, r *http.Request) {
	fmt.Println(r.URL)
	base := parser.Read(*temppath, false)
	contents := parser.Read(*markpath, true)
	output := strings.Replace(base, ContentsKeyword, contents, 1)
	fmt.Fprintf(w, output)
}

func main() {
	// コマンド実行時の引数をパース
	flag.Parse()

	// 実行ファイルのパス
	dir := path.Dir(os.Args[0])
	if *temppath == "static/index.html" {
		*temppath = dir + "/" + *temppath
	}

	// 設定を表示
	fmt.Printf("Markdown path: %s\n", *markpath)
	fmt.Printf("Template path: %s\n", *temppath)
	fmt.Printf("url: http://localhost:%d/\n", *port)

	// workers
	hub := websocket_server.NewHub()
	go hub.Run()
	go worker.Run(hub)

	println(dir + "/static")

	// routing
	http.HandleFunc("/", handler) // ハンドラを登録してウェブページを表示させる
	http.Handle("/static/",
		http.StripPrefix("/static", http.FileServer(http.Dir(dir+"/static")))) // アクセスされたURLから /static 部分を取り除いてハンドリングする
	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		websocket_server.ServeWs(hub, w, r)
	})

	err := http.ListenAndServe(":"+fmt.Sprint(*port), nil)
	if err != nil {
		fmt.Println("ListenAndServe: ", err)
	}
}
