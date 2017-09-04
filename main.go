package main

import (
	"flag"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/ceaser/pac/internal/version"
)

type Page struct {
	Body    []byte
	Message string
}

func (p *Page) save() error {
	return ioutil.WriteFile("pac.js", p.Body, 0600)
}

func loadPage() (*Page, error) {
	body, err := ioutil.ReadFile("pac.js")
	if err != nil {
		return nil, err
	}

	return &Page{Body: body}, nil
}

func viewHandler(w http.ResponseWriter, r *http.Request) {
	p, _ := loadPage()
	w.Header().Set("Content-Type", "application/javascript; charset=utf-8")
	fmt.Fprintf(w, "%s", p.Body)
}

func editHandler(w http.ResponseWriter, r *http.Request) {
	p, err := loadPage()
	if err != nil {
		p = &Page{}
	}
	t, _ := template.ParseFiles("edit.html")
	t.Execute(w, p)
}

func saveHandler(w http.ResponseWriter, r *http.Request) {
	body := r.FormValue("body")
	log.Println(body)
	p := &Page{Body: []byte(body)}
	err := p.save()
	if err != nil {
		p.Message = fmt.Sprintf("Error: %s", err.Error())
		t, _ := template.ParseFiles("edit.html")
		t.Execute(w, p)
	} else {
		http.Redirect(w, r, "/", http.StatusFound)
	}
}

func main() {
	flag.Parse()
	version.ShowVersion()

	log.SetOutput(os.Stdout)

	http.HandleFunc("/", viewHandler)
	http.HandleFunc("/edit/", editHandler)
	http.HandleFunc("/save/", saveHandler)
	http.ListenAndServe(":3000", nil)
}
