/*
Copyright 2017 sycki.
Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at
    http://www.apache.org/licenses/LICENSE-2.0
Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package storage

import (
	"bufio"
	"github.com/sycki/mknote/logger"
	"github.com/sycki/mknote/storage/structs"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"

	"context"
	"github.com/howeyc/fsnotify"
	"github.com/sycki/mknote/cmd/mknote/options"
)

const (
	artTimeFormat = "2006-01-02"
)

type Manager struct {
	ctx           context.Context
	cancel        context.CancelFunc
	config        *options.Config
	lock          *sync.Mutex
	artDir        string                //文章文件所在路径
	latestArticle *structs.Article      //网站首页缓存
	latestIndex   []*structs.ArticleTag //网站索引数据缓存
}

func NewManager(conf *options.Config) *Manager {
	ctx, cancel := context.WithCancel(context.Background())
	return &Manager{
		ctx:    ctx,
		cancel: cancel,
		config: conf,
		artDir: conf.ArticlesDir,
		lock:   &sync.Mutex{},
	}
}

// 从本地文件系统监听是否有文章更新，以更新缓存
func (f *Manager) Start(errCh chan error) {
	logger.Info("starting cache manager ...")
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		errCh <- err
		return
	}

	if err := os.MkdirAll(f.artDir, 0666); err != nil {
		errCh <- err
		return
	}

	//var flags uint32 = syscall.IN_ALL_EVENTS | syscall.IN_CREATE | syscall.IN_DELETE | syscall.IN_MOVE
	var flags uint32 = 0xfff
	err = watcher.AddWatch(f.artDir, flags)
	if err != nil {
		errCh <- err
		return
	}

	//将根文章目录下的所有子目录加入到监听列表
	subDirs, _ := ioutil.ReadDir(f.artDir)
	for _, sub := range subDirs {
		watcher.AddWatch(f.artDir+"/"+sub.Name(), flags)
	}

	go func() {
		for {
			select {
			case <-f.ctx.Done():
				logger.Info("fsmonitor received stop signal")
				watcher.Close()
				return
			case ev := <-watcher.Event:
				//如果根文章目录下创建了新子目录，则加入监听列表，反之从监听列表中删除
				if ev.IsCreate() {
					if f, err := os.Stat(ev.Name); err == nil && f.IsDir() {
						logger.Info("fsmonitor add watch:", ev.Name)
						watcher.AddWatch(ev.Name, flags)
					}
				} else if ev.IsDelete() {
					if f, err := os.Stat(ev.Name); err == nil && f.IsDir() {
						logger.Info("fsmonitor del watch:", ev.Name)
						watcher.RemoveWatch(ev.Name)
					}
				} else if ev.IsRename() {
					if f, err := os.Stat(ev.Name); err == nil && f.IsDir() {
						logger.Info("fsmonitor del watch:", ev.Name)
						watcher.RemoveWatch(ev.Name)
					}
				} else {
					continue
				}
				logger.Info("update cache from fsmonitor event:", ev.String())
				f.latestIndex = nil
				f.latestArticle = nil
			case err := <-watcher.Error:
				logger.Error("fsmonitor read error:", err)
			}
		}
	}()
}

func (f *Manager) Stop() {
	logger.Info("stopping cache manager ...")
	f.cancel()
}

func (f *Manager) GetTitle(uri string) (r string, e error) {
	file, err := os.OpenFile(f.artDir+"/"+uri+".md", os.O_RDONLY, 0666)
	defer file.Close()
	if err != nil {
		logger.Error("failed occur while get article title:", uri)
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
func (f *Manager) GetTags() ([]*structs.ArticleTag, error) {
	if f.latestIndex != nil {
		return f.latestIndex, nil
	}

	subDirInfos, _ := ioutil.ReadDir(f.artDir)
	tagArr := []*structs.ArticleTag{}
	for _, subDirInfo := range subDirInfos {
		if subDirInfo.IsDir() {
			subDir := subDirInfo.Name()
			artInfos, _ := ioutil.ReadDir(f.artDir + "/" + subDir)
			artArr := []*structs.Article{}
			for _, artInfo := range artInfos {
				id := "/" + subDir + "/" + artInfo.Name()[:strings.LastIndex(artInfo.Name(), ".")]
				title, _ := f.GetTitle(id)
				art := &structs.Article{Id: id, Title: title}
				artArr = append(artArr, art)
			}
			tag := &structs.ArticleTag{subDir, artArr}
			tagArr = append(tagArr, tag)
		}
	}

	f.latestIndex = tagArr
	return tagArr, nil
}

func (f *Manager) UpdateMetadata(art *structs.Article) (*structs.Article, error) {
	artFile := f.artDir + "/" + art.Id + ".md"
	fileStr, e := ioutil.ReadFile(artFile)
	if e != nil {
		logger.Error("failed update article meta data:", e)
		return nil, e
	}

	// Divide the article into two parts: body text and metadata
	s := strings.Split(string(fileStr), "\n关于\n---\n")
	art.Content = s[0]
	return f.UpdateArtcile(art)
}

func (f *Manager) UpdateArtcile(art *structs.Article) (*structs.Article, error) {
	artFile := f.artDir + "/" + art.Id + ".md"
	artStr := art.Content
	artStr += "\n关于\n---\n"
	artStr += "\n__作者__：" + art.Author + "\n"
	artStr += "\n__阅读__：" + strconv.Itoa(art.Viewer_count) + "\n"
	artStr += "\n__点赞__：" + strconv.Itoa(art.Like_count) + "\n"
	artStr += "\n__创建__：" + art.Create_date + "\n"

	f.lock.Lock()
	defer f.lock.Unlock()
	err := ioutil.WriteFile(artFile, []byte(artStr), 0666)

	//如果本次更新的文章是缓存中的文章，则更新缓存
	if err == nil && f.latestArticle != nil && f.latestArticle.Id == art.Id {
		f.latestArticle = art
	}
	return art, err
}

func (f *Manager) GetArticle(artID string) (*structs.Article, error) {
	if f.latestArticle != nil && artID == f.latestArticle.Id {
		return f.latestArticle, nil
	}

	artFile := f.artDir + "/" + artID + ".md"

	// load article file data
	fileStr, e := ioutil.ReadFile(artFile)
	if e != nil {
		logger.Error("failed get aritcle:", e)
		return nil, e
	}

	var (
		content, author, create_date, title string
		like_count, viewer_count            int
	)

	// Divide the article into two parts: body text and metadata
	s := strings.Split(string(fileStr), "\n关于\n---\n")

	//这里不用加锁，因为短时间内不会加增浏览数和点赞数，即执行多次写操作也不会有所影响
	if len(s) < 2 {
		file, _ := os.Stat(artFile)
		return f.UpdateMetadata(&structs.Article{
			Id:           artID,
			Author:       f.config.ArticleAuthor,
			Viewer_count: 0,
			Like_count:   0,
			Create_date:  file.ModTime().Format(artTimeFormat),
		})
	}

	scan := bufio.NewScanner(strings.NewReader(s[0]))
	for scan.Scan() {
		line := scan.Text()
		if strings.HasPrefix(line, "# ") {
			title = line[2:]
			break
		}
	}

	content = s[0]

	// read remain string by line
	scan = bufio.NewScanner(strings.NewReader(s[1]))
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
		Id:           artID,
		Title:        title,
		Content:      content,
		Author:       author,
		Viewer_count: viewer_count,
		Like_count:   like_count,
		Create_date:  create_date,
	}, nil
}

func (f *Manager) GetLatestArticleID() (string, error) {
	if f.latestArticle != nil {
		return f.latestArticle.Id, nil
	}

	subDirs, err := ioutil.ReadDir(f.artDir)
	if err != nil {
		logger.Error("Failed get latest article:", err)
		return "", err
	}
	var latestTime time.Time
	var latestArt *structs.Article

	//遍历根文章目录下的所有子目录和文件
	for _, subDir := range subDirs {
		if subDir.IsDir() {
			arts, e1 := ioutil.ReadDir(f.artDir + "/" + subDir.Name())
			if e1 != nil {
				logger.Error("Failed get latest article:", e1)
				return "", e1
			}
			//遍历某子目录下的所有文章
			for _, art := range arts {
				artFile, _ := f.GetArticle("/" + subDir.Name() + "/" + art.Name()[:strings.LastIndex(art.Name(), ".")])
				//获取该文章的创建时间
				//如果文章内没有元数据或没有时间信息，则使用文件的修改时间代替
				createTime, err := time.Parse(artTimeFormat, artFile.Create_date)
				if err != nil {
					logger.Error("Failed get article create date:", err)
					createTime = art.ModTime()
				}

				if latestArt == nil || createTime.After(latestTime) {
					latestTime = createTime
					latestArt = artFile
				}
			}
		} else {
			artFile, _ := f.GetArticle("/" + subDir.Name() + "/" + subDir.Name()[:strings.LastIndex(subDir.Name(), ".")])
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

	f.latestArticle = latestArt
	return latestArt.Id, nil
}
