package view

import (
	"github.com/sycki/mknote/cmd/mknote/options"
	"github.com/sycki/mknote/logger"
	"html/template"
	"io"
	"io/ioutil"
	"net/http"
	"path"
	"strings"
)

type View struct {
	config    *options.Config
	templDir  string
	suffix    string
	templates map[string]*template.Template
}

func NewView(conf *options.Config) *View {
	html := ".html"
	templates := make(map[string]*template.Template)
	fileArr, err := ioutil.ReadDir(conf.HtmlDir)
	check(err)
	var filePath, fileName string
	for _, fileInfo := range fileArr {
		fileName = fileInfo.Name()
		if suffix := path.Ext(fileName); suffix != html {
			continue
		}
		filePath = conf.HtmlDir + "/" + fileName
		t := template.Must(template.ParseFiles(filePath))
		templates[strings.TrimSuffix(fileName, html)] = t
	}

	v := &View{
		config:    conf,
		templDir:  conf.HtmlDir,
		suffix:    html,
		templates: templates,
	}

	return v
}

func check(err error) {
	if err != nil {
		panic(err)
	}
}

func (v *View) RendHTML(w http.ResponseWriter, templ string, model *map[string]interface{}) {
	t, ok := v.templates[templ]
	if !ok {
		logger.Error("Not fount template:", templ)
		return
	}

	err := t.Execute(w, model)
	check(err)
}

func (v *View) SendHTML(w http.ResponseWriter, templ string) {
	fileName := v.templDir + "/" + templ + "/" + v.suffix
	htmlFile, err := ioutil.ReadFile(fileName)
	io.WriteString(w, string(htmlFile))
	check(err)
}
