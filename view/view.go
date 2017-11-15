package view

import (
	"html/template"
	"io"
	"io/ioutil"
	"net/http"
	"path"
	"strings"
)

const (
	TEMPL_DIR = "static/template"
	SUFFIX    = ".html"
)

var templates map[string]*template.Template

func init() {
	templates = make(map[string]*template.Template)
	fileArr, err := ioutil.ReadDir(TEMPL_DIR)
	check(err)
	var filePath, fileName string
	for _, fileInfo := range fileArr {
		fileName = fileInfo.Name()
		if suffix := path.Ext(fileName); suffix != SUFFIX {
			continue
		}
		filePath = TEMPL_DIR + "/" + fileName
		t := template.Must(template.ParseFiles(filePath))
		templates[strings.TrimSuffix(fileName, SUFFIX)] = t
	}
}

func check(err error) {
	if err != nil {
		panic(err)
	}
}

func RenderHTML(w http.ResponseWriter, templ string, model map[string]interface{}) {
	err := templates[templ].Execute(w, model)
	check(err)
}

func SendHTML(w http.ResponseWriter, templ string) {
	htmlFile, err := ioutil.ReadFile(templ + "/" + SUFFIX)
	io.WriteString(w, string(htmlFile))
	check(err)
}
