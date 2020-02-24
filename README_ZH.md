# Marknote
mknote是一个易于使用的博客系统，没有数据库，也无需登录，只需指定你的Markdown文件所在的目录，mknote将自动发现它们，并生成页面。

## 文档
* [中文文档](https://github.com/sycki/mknote/blob/master/README_ZH.md)
* [English doc](https://github.com/sycki/mknote)

## 快速开始
### 安装mknote
进入[发布页面](https://github.com/sycki/mknote/releases)下载最新版本的程序包，然后解压到你喜欢的目录。
```
tar -zxf mknote-<version>.tar.gz
cd mknote-<version>
```

### 启动mknote
通常以https方式启动，这时你需要指定你的证书文件，启动后，它会自动将80端口的请求重定到443端口：
```
bin/mknote \
--tls \
--tls-cert /etc/ssl/cert.pem \
--tls-key /etc/ssl/key.pem
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

## 调试
打开调试功能
```
curl -X POST -H "<your_header_key>: <value>" https://<hostname>/v1/manage/pprof/open
```

使用`go profile`工具进行分析
```
go tool pprof http://<hostname>:8000/debug/pprof/profile
```

关闭调试功能
```
curl -X POST -H "<your_header_key>: <value>" https://<hostname>/v1/manage/pprof/close
```
