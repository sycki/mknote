package blog

import (
	"mknote/config"
	"mknote/structs"
	"net/http"
)

const (
	GET    = "GET"
	POST   = "POST"
	DELETE = "DELETE"
	PUT    = "PUT"
)

var (
	debug   bool
	model   *structs.Model
	htmlDir string
)

func init() {
	debug = true
	model = &structs.Model{make(map[string]interface{})}
	htmlDir = config.Get("html.dir")
}

//func home(w http.ResponseWriter, r *http.Request) {
//	method := r.Method
//	if method == GET {
//		if !isLatestIndex {
//			result, err := database.Index()
//			if err != nil {
//				log.Fatal(err)
//			}
//			index = result
//			isLatestIndex = true
//		}
//		model.Clear()
//		model.Set("index", index)
//		view.SendHTML(w, "home")
//	}
//}

func Article(w http.ResponseWriter, r *http.Request) {
	method := r.Method
	if method == GET {
		//		result, err := database.GetArticle(vars["tag"], vars["en_name"])
		//		if err != nil {
		//			http.Error(w, err.Error(), http.StatusNotFound)
		//			return
		//		}
		//		js, _ := json.Marshal(result)
		//		w.Write(js)
		//		h, err := ioutil.ReadFile(htmlDir + "/article.html")
		//		if err != nil {
		//			http.Error(w, err.Error(), http.StatusNotFound)
		//			return
		//		}
		//		w.Write(h)
		htmlFile := htmlDir + "/article.html"
		http.ServeFile(w, r, htmlFile)
	}
}
