package config

import (
	"bufio"
	"os"
	"regexp"
	"strings"
)

type config struct {
	conf map[string]string
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

// must set keys:
// SYCKIWEB_HOME
// LOG_FILE
func NewConfig() *config {
	file, err := os.OpenFile("config/mknode.conf", os.O_RDONLY, 0666)
	defer file.Close()
	if err != nil {
		panic(err)
	}
	conf := &config{make(map[string]string)}
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
		Set(kv[0], kv[1])
	}
	return conf
}
