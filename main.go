package main

import (
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"regexp"
)

type SvgFile struct {
	Name    string
	SvgCode string
}

type Page struct {
	Path        string
	PageContent []SvgFile
}

var templates = template.Must(template.ParseFiles("index.html"))

func renderPage(path string) (*Page, error) {
	files, err := ioutil.ReadDir("./files")

	if err != nil {
		log.Fatal(err)
	}

	svgFiles := make([]SvgFile, len(files))

	for i, f := range files {
		svgFiles[i].Name = f.Name()

		fileContent, err := ioutil.ReadFile("./files/" + f.Name())

		if err != nil {
			return nil, err
		}

		svgFiles[i].SvgCode = string(fileContent)
	}

	return &Page{Path: path, PageContent: svgFiles}, nil
}

func renderTemplate(w http.ResponseWriter, tmpl string, page *Page) {
	err := templates.ExecuteTemplate(w, tmpl+".html", page)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func makeHandler(fn func(http.ResponseWriter, *http.Request, string)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		m := regexp.MustCompile("^/(files)/([a-zA-Z0-9]+)$").FindStringSubmatch(r.URL.Path)
		if m == nil {
			http.NotFound(w, r)
			return
		}
		fn(w, r, m[2])
	}
}

func explorerHandler(w http.ResponseWriter, r *http.Request, path string) {
	renderedPage, err := renderPage(path)
	if err != nil {
		return
	}
	renderTemplate(w, "index", renderedPage)
}

func main() {
	http.HandleFunc("/files/", makeHandler(explorerHandler))

	log.Fatal(http.ListenAndServe(":8080", nil))
}
