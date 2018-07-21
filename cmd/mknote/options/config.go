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

package options

import (
	"flag"
	"os"
	"path/filepath"
)

var Instance *Config

type Config struct {
	HostName      string
	HomeDir       string
	LogLevel      int
	ArticlesDir   string
	DownloadDir   string
	StaticDir     string
	HtmlDir       string
	IsTls         bool
	TlsCertFile   string
	TlsKeyFile    string
	ArticleAuthor string
	Version       bool
}

func NewDefaultConfig() *Config {
	// setting up default work location
	self, _ := filepath.Abs(os.Args[0])
	binDir := filepath.Dir(self)
	workDir := filepath.Dir(binDir)
	Instance = &Config{
		"",
		workDir,
		1,
		workDir + "/articles",
		workDir + "/f",
		workDir + "/static",
		workDir + "/static/template",
		true,
		workDir + "/conf/fullchain.pem",
		workDir + "/conf/privkey.pem",
		"sycki",
		false,
	}

	return Instance
}

func (c *Config) AddFlags(cmd *flag.FlagSet) {
	cmd.StringVar(&c.HostName, "hostname", c.HostName, "binding hostname")
	cmd.IntVar(&c.LogLevel, "log-level", c.LogLevel, "set log output level, 0...4")
	cmd.StringVar(&c.ArticlesDir, "articles-dir", c.ArticlesDir, "markdown files dir")
	cmd.StringVar(&c.DownloadDir, "uploads-dir", c.DownloadDir, "refence images files dir of all articles")
	cmd.StringVar(&c.StaticDir, "static-dir", c.StaticDir, "html css js etc.")
	cmd.StringVar(&c.HtmlDir, "html-dir", c.HtmlDir, "html template dir")
	cmd.BoolVar(&c.IsTls, "tls", c.IsTls, "specify model https or http")
	cmd.StringVar(&c.TlsCertFile, "tls-cert", c.TlsCertFile, "server cert file for https")
	cmd.StringVar(&c.TlsKeyFile, "tls-key", c.TlsKeyFile, "server key file for https")
	cmd.StringVar(&c.ArticleAuthor, "author", c.ArticleAuthor, "generate author name for article metadata")
	cmd.BoolVar(&c.Version, "version", c.Version, "print version and exit")
}
