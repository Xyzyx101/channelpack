package main

import (
	"fmt"
	"html/template"
	"net/http"
	"path/filepath"
)

// ServePage serves the static web page in /static and deals with templated html
// func ServePage(addr string) {
// 	fs := http.FileServer(http.Dir("static"))
// 	http.Handle("/static/", http.StripPrefix("/static/", fs))
// 	http.HandleFunc("/process", Test)
// 	http.HandleFunc("/", serveTemplate)

// 	http.ListenAndServe(addr, nil)
// 	log.Println("Listening at " + addr + "...")
// }

// func ServIndex(w http.ResponseWriter, r *http.Request) {
// 	lp := filepath.Join("templates", "layout.gohtml")
// 	fp := filepath.Join("templates", "index.gohtml")
// 	tmpl, err := template.ParseFiles(lp, fp)
// 	if err != nil {
// 		fmt.Println(err)
// 		return
// 	}
// 	tmpl.ExecuteTemplate(w, "layout", nil)
// }

// func ServStatic(w http.ResponseWriter, r *http.Request) {
// 	lp := filepath.Join("templates", "layout.gohtml")
// 	fp := filepath.Join("templates", filepath.Clean(r.URL.Path))
// 	tmpl, err := template.ParseFiles(lp, fp)
// 	if err != nil {
// 		fmt.Println(err)
// 		return
// 	}
// 	tmpl.ExecuteTemplate(w, "layout", nil)
// }

// ServeStatic parses templates as serves the static resources and form
func ServeStatic(w http.ResponseWriter, r *http.Request) {
	lp := filepath.Join("templates", "layout.gohtml")
	var fp string
	if r.URL.Path == "/" {
		fp = filepath.Join("templates", "index.gohtml")
	} else {
		fp = filepath.Join("templates", filepath.Clean(r.URL.Path))
	}

	tmpl, err := template.ParseFiles(lp, fp)
	if err != nil {
		fmt.Println(err)
		return
	}
	tmpl.ExecuteTemplate(w, "layout", nil)
}
