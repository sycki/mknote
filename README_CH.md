# Marknote
mknote是一个简单快速的博客系统，没有数据库，也无需登录，只需指定你的Markdown文件所在的目录，mknote将自动发现它们，并生成页面。

## 文档
* [中文文档](https://github.com/sycki/mknote/blob/master/README_CH.md)
* [English doc](https://github.com/sycki/mknote)

## 快速开始
开始之前请确认你已经安装了golang，并配置好了GOPATH环境变量。

### 安装mknote
其中`/usr/local/mknote`指安装路径，且必须是一个不存在的目录
```
go get github.com/sycki/mknote
$GOPATH/src/github.com/sycki/mknote/build.sh install /usr/local/mknote
```

### 确认已经安装好
```
cd /usr/local/mknote && ls
articles  bin  conf  f  static
```

### 启动mknote
通常以https方式启动，这时你需要指定你的证书文件和最终的域名
```
bin/mknote \
--hostname blog.domain.com \
--tls=true \
--tls-cert /etc/ssl/cert.pem \
--tls-key /etc/ssl/key.pem
```

或者以http方式启动它
```
bin/mknote
```

mknote提供了许多有用的选项，用以下命令查看所有选项
```
bin/mknote --help
```

### 编写文章
新建一个文章分类，并编写你的第一篇文章
```
mkdir articles/java
echo '# 第一篇文章' > articles/java/first.md
```

现在你可以在浏览器中访问`http://localhost`

## 文件下载服务
如果你的文章内引用了图片，只需将你的图片文放到`/usr/local/mknote/f/`目录下，该目录可以在启动mknote时指定
```
cp scenery.png /usr/local/mknote/f/
```

然后在文章内引用它
```
![scenery](/f/scenery.png)
```

实际上这个目录中你可以放置任意文件，这样别人就可以在任意地方下载这个文件，这个功能对很多人来说非常实用，在`/usr/local/mknote/f/`目录下可以添加任意多的目录来区分你的文件。

## 引用和参考
* https://github.com/howeyc/fsnotify
* https://github.com/russross/blackfriday
* https://github.com/sindresorhus/github-markdown-css
