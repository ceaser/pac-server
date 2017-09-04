package main

import (
	"flag"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/ceaser/pac-server/internal/version"
)

var (
	templatePath = flag.String("templatepath", "/usr/share/pac-server/tmpl", "Folder where html templates are stored")
	pacFile      = flag.String("pacfile", "/var/spool/pac-server/pac.js", "Location to store PAC file")
	addr         = flag.String("address", ":80", "Address and port to bind to")
)

type Page struct {
	Body    []byte
	Message string
}

func (p *Page) save() error {
	log.Println(*pacFile)
	return ioutil.WriteFile(*pacFile, p.Body, 0600)
}

func (p *Page) editTemplate(w http.ResponseWriter) {
	t, _ := template.ParseFiles(*templatePath + "/edit.html")
	t.Execute(w, p)
}

func loadPage() (*Page, error) {
	log.Println(*pacFile)

	body, err := ioutil.ReadFile(*pacFile)
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
	p.editTemplate(w)
}

func saveHandler(w http.ResponseWriter, r *http.Request) {
	body := r.FormValue("body")
	log.Println(body)
	p := &Page{Body: []byte(body)}
	err := p.save()
	if err != nil {
		p.Message = fmt.Sprintf("Error: %s", err.Error())
		p.editTemplate(w)
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
	http.ListenAndServe(*addr, nil)
}
