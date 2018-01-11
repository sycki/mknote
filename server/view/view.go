package view

import (
	"html/template"
	"io"
	"io/ioutil"
	"mknote/server/ctx"
	"net/http"
	"path"
	"strings"
)

var (
	TEMPL_DIR = ctx.Config.HtmlDir
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

func RendHTML(w http.ResponseWriter, templ string, model *map[string]interface{}) {
	err := templates[templ].Execute(w, *model)
	check(err)
}

func SendHTML(w http.ResponseWriter, templ string) {
	fileName := TEMPL_DIR + "/" + templ + "/" + SUFFIX
	htmlFile, err := ioutil.ReadFile(fileName)
	io.WriteString(w, string(htmlFile))
	check(err)
}
