package main

import (
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"regexp"
)

type Folder struct {
	Name string
}

type SvgFile struct {
	Content template.HTML
	Name    string
}

type Page struct {
	Folders  []Folder
	Path     string
	SvgFiles []SvgFile
}

var templates = template.Must(template.ParseFiles("index.html"))

func renderPage(path string) (*Page, error) {
	files, err := ioutil.ReadDir("./files")

	if err != nil {
		log.Fatal(err)
	}

	folders := []Folder{}
	svgFiles := []SvgFile{}

	for _, file := range files {
		if file.IsDir() {
			folders = append(folders, Folder{Name: file.Name()})
		} else {
			content, err := ioutil.ReadFile("./files/" + file.Name())

			if err != nil {
				return nil, err
			}

			svgFiles = append(svgFiles, SvgFile{Content: template.HTML(content), Name: file.Name()})
		}
	}

	return &Page{Folders: folders, Path: path, SvgFiles: svgFiles}, nil
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

func heheHandler(w http.ResponseWriter, r *http.Request, path string) {
	http.Error(w, "error ;) hhe", http.StatusInternalServerError)
}

func main() {
	http.HandleFunc("/files/", makeHandler(explorerHandler))
	http.Handle("/hehe", makeHandler(heheHandler))
	http.Handle("/", http.StripPrefix("/", http.FileServer(http.Dir(""))))

	log.Fatal(http.ListenAndServe(":8080", nil))
}
