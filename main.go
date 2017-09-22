package main

import (
	"context"
	"crypto/md5"
	"encoding/hex"
	"flag"
	"fmt"
	"html/template"
	"io/ioutil"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"time"

	"github.com/ceaser/pac-server/internal/version"
	"github.com/ceaser/pac-server/logging"
)

var (
	templatePath = flag.String("templatepath", "/usr/share/pac-server/tmpl", "Folder where html templates are stored")
	pacFile      = flag.String("pacfile", "/var/spool/pac-server/pac.js", "Location to store PAC file")
	addr         = flag.String("address", ":80", "Address and port to bind to")
	maxAge       = flag.String("maxage", "31536000", "Cache Control max-age")

	lastEtag string
)

type Page struct {
	Body    []byte
	Message string
	Etag    string
}

func (p *Page) save() error {
	return ioutil.WriteFile(*pacFile, p.Body, 0600)
}

func (p *Page) editTemplate(w http.ResponseWriter) {
	t, _ := template.ParseFiles(*templatePath + "/edit.html")
	t.Execute(w, p)
}

func loadPage() (*Page, error) {
	body, err := ioutil.ReadFile(*pacFile)
	if err != nil {
		return nil, err
	}

	page := &Page{Body: body}

	// Generate etag
	hasher := md5.New()
	hasher.Write([]byte(page.Body))
	page.Etag = hex.EncodeToString(hasher.Sum(nil))

	return page, nil
}

func wpadHandler(w http.ResponseWriter, r *http.Request) {
	// TODO: Refactor duplication between wpadHandler and viewHandler
	p, _ := loadPage()
	e := p.Etag

	// Caching
	w.Header().Set("Etag", e)
	w.Header().Set("Cache-Control", "max-age="+*maxAge)
	if match := r.Header.Get("If-None-Match"); match != "" {
		if strings.Contains(match, e) {
			w.WriteHeader(http.StatusNotModified)
			return
		}
	}

	w.Header().Set("Content-Type", "application/x-ns-proxy-autoconfig; charset=utf-8")
	fmt.Fprintf(w, "%s", p.Body)
}

func viewHandler(w http.ResponseWriter, r *http.Request) {
	p, _ := loadPage()
	e := p.Etag

	// Caching
	w.Header().Set("Etag", e)
	w.Header().Set("Cache-Control", "max-age="+*maxAge)
	if match := r.Header.Get("If-None-Match"); match != "" {
		if strings.Contains(match, e) {
			w.WriteHeader(http.StatusNotModified)
			return
		}
	}

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
	p := &Page{Body: []byte(body)}
	err := p.save()
	if err != nil {
		p.Message = fmt.Sprintf("Error: %s", err.Error())
		p.editTemplate(w)
	} else {
		http.Redirect(w, r, "/", http.StatusFound)
	}
}

func missingHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotFound)
	fmt.Fprint(w, "404")
}

func main() {
	flag.Parse()
	version.ShowVersion()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)

	mux := http.DefaultServeMux
	mux.HandleFunc("/", viewHandler)
	mux.HandleFunc("/wpad.dat", wpadHandler)
	mux.HandleFunc("/edit/", editHandler)
	mux.HandleFunc("/save/", saveHandler)
	mux.HandleFunc("/favicon.ico", missingHandler)
	loggingHandler := logging.NewApacheLoggingHandler(mux, os.Stdout)
	server := &http.Server{
		Addr:    *addr,
		Handler: loggingHandler,
	}

	go func() {
		server.ListenAndServe()
	}()
	<-stop
	ctx, _ := context.WithTimeout(context.Background(), 9*time.Second)
	server.Shutdown(ctx)
}
