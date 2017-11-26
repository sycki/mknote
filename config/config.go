package config

import (
	"bufio"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

type config struct {
	conf map[string]string
}

func (c *config) addDefault(k string, v string) {
	println("add default configuration", k, "=>", v)
	c.conf[k] = v
}

func Set(k string, v string) {
	c.conf[k] = v
}

func Get(k string) string {
	v, _ := c.conf[k]
	return v
}

func GetOr(k string, d string) string {
	v, ok := c.conf[k]
	if ok {
		return v
	} else {
		return d
	}
}

var (
	c *config
)

func init() {
	c = NewConfig()
}

// must exists keys:
// MKNOTE_HOME
// log.file
func NewConfig() *config {
	// create config object and load default properties
	conf := &config{make(map[string]string)}

	// setting up default work location
	self, _ := filepath.Abs(os.Args[0])
	workDir := filepath.Dir(self)
	conf.addDefault("MKNOTE_HOME", workDir)

	// setting up to default
	conf.addDefault("log.file", workDir+"/log/mknote.log")
	conf.addDefault("articles.dir", workDir+"/articles")
	conf.addDefault("html.dir", workDir+"/static/template")
	conf.addDefault("article.default.author", "sycki")

	// load config file, exit the program if config file not exists
	file, err := os.OpenFile(workDir+"/conf/mknote.conf", os.O_RDONLY, 0666)
	defer file.Close()
	if err != nil {
		panic(err)
	}

	scan := bufio.NewScanner(file)
	r := regexp.MustCompile(`[ \t]+`)
	for scan.Scan() {
		line := scan.Text()
		kv := r.Split(line, 2)
		if len(kv) < 2 || strings.HasPrefix(kv[0], "#") {
			continue
		}
		conf.addDefault(kv[0], kv[1])
	}
	return conf
}
