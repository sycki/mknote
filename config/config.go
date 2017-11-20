package config

import (
	"bufio"
	"os"
	"regexp"
	"strings"
)

type Config struct {
	conf map[string]string
}

func (c *Config) Set(k string, v string) {
	c.conf[k] = v
}

func (c *Config) Get(k string) string {
	v, _ := c.conf[k]
	return v
}

func (c *Config) GetOr(k string, d string) string {
	v, ok := c.conf[k]
	if ok {
		return v
	} else {
		return d
	}
}

var (
	Conf *Config
)

func init() {
	Conf = NewConfig()
}

// must set keys:
// SYCKIWEB_HOME
func NewConfig() *Config {
	file, err := os.OpenFile("config/syckiweb.conf", os.O_RDONLY, 0666)
	defer file.Close()
	if err != nil {
		panic(err)
	}
	conf := &Config{make(map[string]string)}
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
		conf.Set(kv[0], kv[1])
	}
	return conf
}
