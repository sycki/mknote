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

package ctx

import (
	"flag"
	"os"
	"path/filepath"
)

var (
	version = "mknote-v2.2.3"
	Config  *config
)

type config struct {
	HomeDir       string
	LogFile       string
	LogLevel      int
	ArticlesDir   string
	UploadsDir    string
	StaticDir     string
	HtmlDir       string
	TlsCertFile   string
	TlsKeyFile    string
	ArticleAuthor string
}

func init() {
	// setting up default work location
	self, _ := filepath.Abs(os.Args[0])
	workDir := filepath.Dir(self)

	Config = &config{
		workDir,
		workDir + "/log/mknote.log",
		1,
		workDir + "/articles",
		workDir + "/uploads",
		workDir + "/static",
		workDir + "/static/template",
		workDir + "/conf/fullchain.pem",
		workDir + "/conf/privkey.pem",
		"sycki",
	}

	for _, arg := range os.Args {
		if arg == "--version" {
			println(version)
			os.Exit(0)
		}
	}

	flag.StringVar(&Config.LogFile, "log-file", Config.LogFile, "set log output file")
	flag.IntVar(&Config.LogLevel, "log-level", Config.LogLevel, "set log output level, 0...4")
	flag.StringVar(&Config.ArticlesDir, "articles-dir", Config.ArticlesDir, "markdown files dir")
	flag.StringVar(&Config.UploadsDir, "uploads-dir", Config.UploadsDir, "refence images files dir of all articles")
	flag.StringVar(&Config.StaticDir, "static-dir", Config.StaticDir, "css js etc.")
	flag.StringVar(&Config.HtmlDir, "html-dir", Config.HtmlDir, "html template dir")
	flag.StringVar(&Config.TlsCertFile, "tls-cert", Config.TlsCertFile, "server cert file for https")
	flag.StringVar(&Config.TlsKeyFile, "tls-key", Config.TlsKeyFile, "server key file for https")
	flag.StringVar(&Config.ArticleAuthor, "author", Config.ArticleAuthor, "generate author name for article metadata")

	flag.Parse()

	for _, arg := range os.Args {
		if arg == "--help" {
			os.Exit(0)
		}
	}

}
