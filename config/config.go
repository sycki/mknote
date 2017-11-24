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

	// setting up default log file
	sep := string(os.PathSeparator)
	logFile := workDir + sep + "log" + sep + "mknote.log"
	conf.addDefault("log.file", logFile)

	// load config file, exit the program if config file not exists
	file, err := os.OpenFile("config/mknote.conf", os.O_RDONLY, 0666)
	defer file.Close()
	if err != nil {
		panic(err)
	}

	in := bufio.NewReader(file)
	for {
		line, err := in.ReadString('\n')
		if err != nil {
			break
		}
		kv := regexp.MustCompilePOSIX("\\s+").Split(line, 2)
		if len(kv) < 2 || strings.HasPrefix(kv[0], "#") {
			continue
		}
		conf.addDefault(kv[0], kv[1])
	}
	return conf
}
