package httpserver

import (
	"html/template"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/parth105/simple-http/internal/wikipage"
)

// Consolidate all the template parsing
// template.Must panics if parsing returns non-nil error
var templates *template.Template

// Called upon by each of the handlers to respond (render) to the client/requester
// with a page.
func renderTemplate(w http.ResponseWriter, t string, p *wikipage.Page) {
	template, err := template.ParseFiles(t)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	err = template.Execute(w, p)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

// Default handler for the root "/" URL
func handler(w http.ResponseWriter, r *http.Request) {
	files, _ := os.ReadDir(".")
	var pages []string

	for _, file := range files {
		if strings.Contains(file.Name(), ".page") {
			pages = append(pages, strings.Replace(file.Name(), ".page", "", -1))
		}
	}
	t := template.Must(template.ParseFiles("web/welcome.html"))
	t.Execute(w, pages)
}

// Handler for viewing a wiki page
func viewHandler(w http.ResponseWriter, r *http.Request) {
	title := r.URL.Path[len("/view/"):]
	p, err := wikipage.LoadPage(title)
	if err != nil {
		http.Redirect(w, r, "/edit/"+title, http.StatusFound)
		return
	}
	renderTemplate(w, "web/view.html", p)
}

// Handler for editing an existing page or creating a new one
func editHandler(w http.ResponseWriter, r *http.Request) {
	title := r.URL.Path[len("/edit/"):]
	p, err := wikipage.LoadPage(title)
	if err != nil {
		p = &wikipage.Page{Title: title}
	}
	renderTemplate(w, "web/edit.html", p)
}

// Handler for saving a wiki page to the disk
func saveHandler(w http.ResponseWriter, r *http.Request) {
	title := r.URL.Path[len("/save/"):]
	content := r.FormValue("body")
	p := &wikipage.Page{Title: title, Body: []byte(content)}
	err := p.Save()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	http.Redirect(w, r, "/view/"+title, http.StatusFound)
}

// Starts the wiki page server with handlers
func StartServer(port string) {
	// The static file paths are relative to the main.go
	// and the name of the templates are the filenames (including the web path)
	templates = template.Must(template.ParseFiles("web/welcome.html", "web/view.html", "web/edit.html"))
	http.HandleFunc("/", handler)
	http.HandleFunc("/view/", viewHandler)
	http.HandleFunc("/edit/", editHandler)
	http.HandleFunc("/save/", saveHandler)
	//http.HandleFunc("/list/", viewHandler)
	if port == "" {
		port = "8089"
	}
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
