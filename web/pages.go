package web

import (
	"html/template"
	"net/http"
	"os"
)

func runLocalhost() {
	// place pages here, using http.HandleFunc
	http.HandleFunc("/home/", homeHandler)
	http.HandleFunc("/dashboard/", dashboardHandler)
	http.HandleFunc("/issues/", issuesHandler)
	http.HandleFunc("/settings/", settingsHandler)

	// make sure http is able to use static files (i.e. css and js)
	http.Handle("/Resources/App/", http.StripPrefix("/Resources/App/",
		http.FileServer(http.Dir("./Resources/App"))))
}

type Page struct {
	Title string
	Body  []byte
}

func loadPage(title string) (*Page, error) {
	filename := title + ".txt"
	body, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	return &Page{Title: title, Body: body}, nil
}

func renderTemplate(w http.ResponseWriter, tmpl string, p *Page) {
	t, _ := template.ParseFiles(tmpl + ".html")
	t.Execute(w, p)
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	title := r.URL.Path[len("/home/"):]
	p, err := loadPage(title)
	if err != nil {
		p = &Page{Title: title}
	}
	renderTemplate(w, "Resources/App/Html/home", p)
}

func dashboardHandler(w http.ResponseWriter, r *http.Request) {
	title := r.URL.Path[len("/dashboard/"):]
	p, err := loadPage(title)
	if err != nil {
		p = &Page{Title: title}
	}
	renderTemplate(w, "Resources/App/Html/dashboard", p)
}

func issuesHandler(w http.ResponseWriter, r *http.Request) {
	title := r.URL.Path[len("/issues/"):]
	p, err := loadPage(title)
	if err != nil {
		p = &Page{Title: title}
	}
	renderTemplate(w, "Resources/App/Html/issues", p)
}

func settingsHandler(w http.ResponseWriter, r *http.Request) {
	title := r.URL.Path[len("/settings/"):]
	p, err := loadPage(title)
	if err != nil {
		p = &Page{Title: title}
	}
	renderTemplate(w, "Resources/App/Html/settings", p)
}
