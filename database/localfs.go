package database

import (
	"bufio"
	"errors"
	"html/template"
	"io/ioutil"
	"log"
	"strconv"
	"strings"
	"sycki/config"
	"sycki/structs"
)

var (
	Templates map[string]*template.Template
)

func init() {
	Templates = make(map[string]*template.Template)
	log.Println("INFO init jsonDAO complete.")
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
	artDir := config.Get("SYCKIWEB_HOME") + "/articles"
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

func GetArticle(tag string, en_name string) (*structs.Article, error) {
	artPath := config.Get("SYCKIWEB_HOME") + "/articles" + "/" + tag + "/" + en_name + ".md"
	//	artFile,e := os.OpenFile(artPath, os.O_RDONLY, 0666)
	//	defer artFile.Close()
	//	if e != nil {
	//		return nil, e
	//	}
	var (
		content, author, create_date string
		like_count, viewer_count     int
	)
	fileStr, e := ioutil.ReadFile(artPath)
	if e != nil {
		return nil, e
	}

	s := strings.Split(string(fileStr), "\n关于\n---\n")
	if len(s) < 2 {
		return nil, errors.New("Error article file is found, but not found metadata.")
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
