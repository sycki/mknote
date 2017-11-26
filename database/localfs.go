package database

import (
	"bufio"
	"html/template"
	"io/ioutil"
	"mknote/config"
	"mknote/log"
	"mknote/structs"
	"strconv"
	"strings"
	"sync"
	"time"
)

var (
	Templates map[string]*template.Template
	artDir    string
)

func init() {
	Templates = make(map[string]*template.Template)
	artDir = config.Get("articles.dir")
	log.Info("init database complete.")
}

/*
构造一个首页文章导航json
[
	{tag1:[art1,art2,...]},
	{tag2:[art3,art4,...]},
	...
]
*/
func GetTags() ([]*structs.ArticleTag, error) {
	subDirInfos, _ := ioutil.ReadDir(artDir)
	tagArr := []*structs.ArticleTag{}
	for _, subDirInfo := range subDirInfos {
		if subDirInfo.IsDir() {
			subDir := subDirInfo.Name()
			artInfos, _ := ioutil.ReadDir(artDir + "/" + subDir)
			artArr := []string{}
			for _, artInfo := range artInfos {
				artArr = append(artArr, artInfo.Name())
			}
			tag := &structs.ArticleTag{subDir, artArr}
			tagArr = append(tagArr, tag)
		}
	}

	return tagArr, nil
}

var l sync.Mutex

func UpdateMetadata(art *structs.Article) (*structs.Article, error) {
	artID := art.ID
	artFile := artDir + artID + ".md"
	fileStr, e := ioutil.ReadFile(artFile)
	if e != nil {
		log.Error("failed update article meta data:", e)
		return nil, e
	}

	// Divide the article into two parts: body text and metadata
	s := strings.Split(string(fileStr), "\n关于\n---\n")
	art.Content = s[0]
	return UpdateArtcile(art)
}

func UpdateArtcile(art *structs.Article) (*structs.Article, error) {
	artID := art.ID
	artFile := artDir + artID + ".md"
	artStr := art.Content
	artStr += "\n关于\n---\n"
	artStr += "\n__作者__：" + art.Author + "\n"
	artStr += "\n__阅读__：" + strconv.Itoa(art.Viewer_count) + "\n"
	artStr += "\n__点赞__：" + strconv.Itoa(art.Like_count) + "\n"
	artStr += "\n__创建__：" + art.Create_date + "\n"

	l.Lock()
	defer l.Unlock()
	err := ioutil.WriteFile(artFile, []byte(artStr), 0666)

	return art, err
}

func GetArticle(artID string) (*structs.Article, error) {
	artFile := artDir + artID + ".md"

	// load article file data
	fileStr, e := ioutil.ReadFile(artFile)
	if e != nil {
		log.Error("failed get aritcle:", e)
		return nil, e
	}

	var (
		content, author, create_date string
		like_count, viewer_count     int
	)

	// Divide the article into two parts: body text and metadata
	s := strings.Split(string(fileStr), "\n关于\n---\n")

	//这里不用加锁，因为短时间内不会加增浏览数和点赞数，即执行多次写操作也不会有所影响
	if len(s) < 2 {
		return UpdateMetadata(&structs.Article{
			ID:           artID,
			Author:       config.Get("article.default.author"),
			Viewer_count: 0,
			Like_count:   0,
			Create_date:  time.Now().Format("2017-01-02"),
		})
	}

	content = s[0]

	// read remain string by line
	scan := bufio.NewScanner(strings.NewReader(s[1]))
	for scan.Scan() {
		line := scan.Text()
		kv := strings.Split(line, "：")
		if len(kv) < 2 {
			continue
		}
		key := kv[0]
		if key == "__作者__" {
			author = kv[1]
		} else if key == "__阅读__" {
			i, err := strconv.Atoi(kv[1])
			if err != nil {
				viewer_count = 0
			} else {
				viewer_count = i
			}
		} else if key == "__点赞__" {
			i, err := strconv.Atoi(kv[1])
			if err != nil {
				like_count = 0
			} else {
				like_count = i
			}
		} else if key == "__创建__" {
			create_date = kv[1]
		}
	}

	return &structs.Article{
		Content:      content,
		Author:       author,
		Viewer_count: viewer_count,
		Like_count:   like_count,
		Create_date:  create_date,
	}, nil
}

//func GetArticle(articleName string) (resutl string, err error) {

//}
