package main

import (
	"fmt"
	"net/http"
	"strings"
	"flag"

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
	fmt.Println(base)
	output := strings.Replace(base, ContentsKeyword, contents, 1)
	fmt.Println(output)
	fmt.Fprintf(w, output)
}

func main() {
	flag.Parse()

	fmt.Printf("Markdown path: %s\n", *markpath)
	fmt.Printf("Template path: %s\n", *temppath)
	fmt.Printf("url: http://localhost:%d/\n", *port)

	// workers
	hub := websocket_server.NewHub()
	go hub.Run()
	go worker.Run(hub)

	// routing
	http.HandleFunc("/", handler) // ハンドラを登録してウェブページを表示させる
	http.Handle("/static/", http.StripPrefix("/static", http.FileServer(http.Dir("static")))) // アクセスされたURLから /static 部分を取り除いてハンドリングする
	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		websocket_server.ServeWs(hub, w, r)
	})

	err := http.ListenAndServe(":" + fmt.Sprint(*port), nil)
	if err != nil {
		fmt.Println("ListenAndServe: ", err)
	}
}
