package options

import (
	"flag"
	"os"
	"path/filepath"
)

var Instance *Config

type Config struct {
	HttpAddrPort   string
	TlsAddrPort    string
	HomeDir        string
	LogLevel       int
	ArticlesDir    string
	DownloadDir    string
	HtmlDir        string
	IsTls          bool
	TlsCertFile    string
	TlsKeyFile     string
	IsRedirectHttp bool
	ArticleAuthor  string
	Version        bool
	DebugPort      string
}

func NewDefaultConfig() *Config {
	// setting up default work location
	self, _ := filepath.Abs(os.Args[0])
	binDir := filepath.Dir(self)
	workDir := filepath.Dir(binDir)
	Instance = &Config{
		":80",
		":443",
		workDir,
		1,
		workDir + "/articles",
		workDir + "/f",
		workDir + "/html",
		false,
		workDir + "/conf/fullchain.pem",
		workDir + "/conf/privkey.pem",
		false,
		"sycki",
		false,
		"8000",
	}

	return Instance
}

// AddFlags add options for command
func (c *Config) AddFlags(cmd *flag.FlagSet) {
	cmd.StringVar(&c.HttpAddrPort, "port-http", c.HttpAddrPort, "http server listen port")
	cmd.StringVar(&c.TlsAddrPort, "port-tls", c.TlsAddrPort, "https server listen port")
	cmd.IntVar(&c.LogLevel, "log-level", c.LogLevel, "set log output level, 0...4")
	cmd.StringVar(&c.ArticlesDir, "articles-dir", c.ArticlesDir, "markdown files dir")
	cmd.StringVar(&c.DownloadDir, "download-dir", c.DownloadDir, "file server directory of all assets")
	cmd.StringVar(&c. HtmlDir, "html-dir", c.HtmlDir, "html template dir")
	cmd.BoolVar(&c.IsTls, "tls", c.IsTls, "specify https model instead of http")
	cmd.StringVar(&c.TlsCertFile, "tls-cert", c.TlsCertFile, "server cert file for https")
	cmd.StringVar(&c.TlsKeyFile, "tls-key", c.TlsKeyFile, "server key file for https")
	cmd.BoolVar(&c.IsRedirectHttp, "redirect-http", c.IsRedirectHttp, "http request redirect to tls port")
	cmd.StringVar(&c.ArticleAuthor, "author", c.ArticleAuthor, "generate author name for article metadata")
	cmd.BoolVar(&c.Version, "version", c.Version, "print version and exit")
	cmd.StringVar(&c.DebugPort, "debug-port", c.DebugPort, "used pprofile server port")
}
