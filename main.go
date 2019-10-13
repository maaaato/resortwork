package main

import (
	"flag"
	"html/template"
	"log"
	"net/http"
	"path/filepath"
	"sync"
	"time"
	"os"
)

type templateHandler struct {
	once     sync.Once
	filename string
	templ    *template.Template
}

type trapHandler struct {
	once     sync.Once
	filename string
	templ    *template.Template
}


func (t *templateHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	t.once.Do(func() {
		t.templ =
			template.Must(template.ParseFiles(filepath.Join("templates",
				t.filename)))
	})
	data := map[string]interface{}{}
	t.templ.Execute(w, data)
}


func (t *trapHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	t.once.Do(func() {
		t.templ =
			template.Must(template.ParseFiles(filepath.Join("templates",
				t.filename)))
	})
	const layout = "2006-01-02 15:04:05"
	jst, _ := time.LoadLocation("Asia/Tokyo")
	access := time.Now().In(jst)
	data := map[string]interface{}{
		"Host": r.Host,
		"IP": r.RemoteAddr,
		"UserAgent": r.Header.Get("User-Agent"),
		"Time": access.Format(layout),
	}
	t.templ.Execute(w, data)
}

func main() {
	addr := os.Getenv("PORT")
	flag.Parse()

	http.Handle("/", &templateHandler{filename: "login.html"})
	http.Handle("/trap", &trapHandler{filename: "trap.html"})
	if err := http.ListenAndServe(":"+addr, nil); err != nil {
		log.Fatal("ListenAndServe:", err)
	}
}