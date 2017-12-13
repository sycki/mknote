package persistent

import (
	"bufio"
	"io/ioutil"
	"mknote/server/ctx"
	"mknote/server/persistent/structs"
	"os"
	"strconv"
	"strings"
	"sync"
	"syscall"
	"time"

	"github.com/howeyc/fsnotify"
)

const (
	artTimeFormat = "2006-01-02"
)

var (
	artDir string
	//网站首页缓存
	latestArticle *structs.Article
	//网站索引数据缓存
	latestIndex []*structs.ArticleTag
	l           = &sync.Mutex{}
)

func init() {
	artDir = ctx.Get("articles.dir")
	go fsmonitor()
	ctx.Info("initialize data base complete")
}

// 监听本地文件系统是否有新的文章，以更新缓存
func fsmonitor() {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		ctx.Fatal(err)
	}

	sigs := make(chan string, 1)
	ctx.RegistryStoper(sigs)

	//var flags uint32 = syscall.IN_CREATE | syscall.IN_DELETE | syscall.IN_MOVE
	var flags uint32 = syscall.IN_ALL_EVENTS
	err = watcher.AddWatch(artDir, flags)
	if err != nil {
		ctx.Fatal(err)
	}

	//将根文章目录下的所有子目录加入到监听列表
	subDirs, _ := ioutil.ReadDir(artDir)
	for _, sub := range subDirs {
		watcher.AddWatch(artDir+"/"+sub.Name(), flags)
	}

	for {
		select {
		case sig := <-sigs:
			ctx.Info("fsmonitor read sigs:", sig)
			break
		case ev := <-watcher.Event:
			//如果根文章目录下创建了新子目录，则加入监听列表，反之删除
			if ev.IsCreate() {
				if f, err := os.Stat(ev.Name); err == nil && f.IsDir() {
					ctx.Info("fsmonitor add watch:", ev.Name)
					watcher.AddWatch(ev.Name, flags)
				}
			} else if ev.IsDelete() {
				if f, err := os.Stat(ev.Name); err == nil && f.IsDir() {
					ctx.Info("fsmonitor del watch:", ev.Name)
					watcher.RemoveWatch(ev.Name)
				}
			} else if ev.IsRename() {
				if f, err := os.Stat(ev.Name); err == nil && f.IsDir() {
					ctx.Info("fsmonitor del watch:", ev.Name)
					watcher.RemoveWatch(ev.Name)
				}
			} else {
				continue
			}
			ctx.Info("fsmonitor read event:", ev.String())
			latestIndex = nil
			latestArticle = nil
		case err := <-watcher.Error:
			ctx.Error("fsmonitor read error:", err)
		}
	}

	close(sigs)
	watcher.Close()

}

func GetTitle(uri string) (r string, e error) {
	file, err := os.OpenFile(artDir+uri+".md", os.O_RDONLY, 0666)
	defer file.Close()
	if err != nil {
		ctx.Error("failed occur while get article title:", uri)
		return "", err
	}
	scan := bufio.NewScanner(file)
	for scan.Scan() {
		line := scan.Text()
		if strings.HasPrefix(line, "# ") {
			r = line[2:]
			break
		}
	}
	return
}

/*
构造一个首页文章导航json
[{"Name":"linux","Articles":[{"ID":"/linux/linux-code1","Title":"mknote","En_name":"","Content":"","Author":"","Like_count":0,"Viewer_count":0,"Create_date":""},{"ID":"/linux/linux-code2","Title":"Linux 管道","En_name":"","Content":"","Author":"","Like_count":0,"Viewer_count":0,"Create_date":""}]},{"Name":"mknote","Articles":[{"ID":"/mknote/README","Title":"mknote","En_name":"","Content":"","Author":"","Like_count":0,"Viewer_count":0,"Create_date":""}]}]
*/
func GetTags() ([]*structs.ArticleTag, error) {
	if latestIndex != nil {
		return latestIndex, nil
	}

	subDirInfos, _ := ioutil.ReadDir(artDir)
	tagArr := []*structs.ArticleTag{}
	for _, subDirInfo := range subDirInfos {
		if subDirInfo.IsDir() {
			subDir := subDirInfo.Name()
			artInfos, _ := ioutil.ReadDir(artDir + "/" + subDir)
			artArr := []*structs.Article{}
			for _, artInfo := range artInfos {
				id := "/" + subDir + "/" + artInfo.Name()[:strings.LastIndex(artInfo.Name(), ".")]
				title, _ := GetTitle(id)
				art := &structs.Article{ID: id, Title: title}
				artArr = append(artArr, art)
			}
			tag := &structs.ArticleTag{subDir, artArr}
			tagArr = append(tagArr, tag)
		}
	}

	latestIndex = tagArr
	return tagArr, nil
}

func UpdateMetadata(art *structs.Article) (*structs.Article, error) {
	artFile := artDir + art.ID + ".md"
	fileStr, e := ioutil.ReadFile(artFile)
	if e != nil {
		ctx.Error("failed update article meta data:", e)
		return nil, e
	}

	// Divide the article into two parts: body text and metadata
	s := strings.Split(string(fileStr), "\n关于\n---\n")
	art.Content = s[0]
	return UpdateArtcile(art)
}

func UpdateArtcile(art *structs.Article) (*structs.Article, error) {
	artFile := artDir + art.ID + ".md"
	artStr := art.Content
	artStr += "\n关于\n---\n"
	artStr += "\n__作者__：" + art.Author + "\n"
	artStr += "\n__阅读__：" + strconv.Itoa(art.Viewer_count) + "\n"
	artStr += "\n__点赞__：" + strconv.Itoa(art.Like_count) + "\n"
	artStr += "\n__创建__：" + art.Create_date + "\n"

	l.Lock()
	defer l.Unlock()
	err := ioutil.WriteFile(artFile, []byte(artStr), 0666)

	//如果本次更新的文章是缓存中的文章，则更新缓存
	if err != nil && latestArticle != nil && latestArticle.ID == art.ID {
		latestArticle = art
	}
	return art, err
}

func GetArticle(artID string) (*structs.Article, error) {
	if latestArticle != nil && artID == latestArticle.ID {
		return latestArticle, nil
	}

	artFile := artDir + artID + ".md"

	// load article file data
	fileStr, e := ioutil.ReadFile(artFile)
	if e != nil {
		ctx.Error("failed get aritcle:", e)
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
		f, _ := os.Stat(artFile)
		return UpdateMetadata(&structs.Article{
			ID:           artID,
			Author:       ctx.Get("article.default.author"),
			Viewer_count: 0,
			Like_count:   0,
			Create_date:  f.ModTime().Format(artTimeFormat),
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
		ID:           artID,
		Content:      content,
		Author:       author,
		Viewer_count: viewer_count,
		Like_count:   like_count,
		Create_date:  create_date,
	}, nil
}

func GetLatestArticleID() (string, error) {
	if latestArticle != nil {
		return latestArticle.ID, nil
	}

	subDirs, err := ioutil.ReadDir(artDir)
	if err != nil {
		ctx.Error("Failed get latest article:", err)
		return "", err
	}
	var latestTime time.Time
	var latestArt *structs.Article

	//遍历根文章目录下的所有子目录和文件
	for _, subDir := range subDirs {
		if subDir.IsDir() {
			arts, e1 := ioutil.ReadDir(artDir + "/" + subDir.Name())
			if e1 != nil {
				ctx.Error("Failed get latest article:", e1)
				return "", e1
			}
			//遍历某子目录下的所有文章
			for _, art := range arts {
				artFile, _ := GetArticle("/" + subDir.Name() + "/" + art.Name()[:strings.LastIndex(art.Name(), ".")])
				//获取该文章的创建时间
				//如果文章内没有元数据或没有时间信息，则使用文件的修改时间代替
				createTime, err := time.Parse(artTimeFormat, artFile.Create_date)
				if err != nil {
					ctx.Error("Failed get article create date:", err)
					createTime = art.ModTime()
				}

				if latestArt == nil || createTime.After(latestTime) {
					latestTime = createTime
					latestArt = artFile
				}
			}
		} else {
			artFile, _ := GetArticle("/" + subDir.Name() + "/" + subDir.Name()[:strings.LastIndex(subDir.Name(), ".")])
			createTime, err := time.Parse(artTimeFormat, artFile.Create_date)
			if err != nil {
				createTime = subDir.ModTime()
			}

			if latestArt == nil || createTime.After(latestTime) {
				latestTime = createTime
				latestArt = artFile
			}
		}
	}

	latestArticle = latestArt
	return latestArt.ID, nil
}
