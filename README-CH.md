# Marknote
mknote是一个简单快速的博客系统，没有数据库，也无需登录，只需指定你的Markdown文件所在的目录，mknote将自动发现它们，并生成页面。

## 文档
* [中文文档](https://github.com/sycki/mknote/blob/master/README-CH.md)
* [English doc](https://github.com/sycki/mknote)

## 开发指南
### 克隆项目
你可以先Fork它，也可以直接克隆它。
```
cd $GOPATH/src
git clone https://github.com/sycki/mknote.git
```

### 开发
导入你的IDE中并开发你想要的功能。

### 编译和运行
```
cd $GOPATH/src/mknote
go build -v mknote
sudo ./mknote
```

## 使用指南
### 下载
你可以直接去发布页面下载已经编译好的二进制程序。
[https://github.com/sycki/mknote/releases](https://github.com/sycki/mknote/releases)

### 安装
```
mkdir /usr/local/mknote/
tar -xf mknote-v2.2.tar -C /usr/local/mknote/
```

### 启动
```
/usr/local/sycki-mknote/mknote &
```

### TLS方式启动
如果你想以https的方式启动它，请在启动时使用`--tls-cert`和`--tls-key`两个选项指定你的证书和私钥文件。
```
/usr/local/sycki-mknote/mknote \
--tls-cert /etc/ssl/cert.pem \
--tls-key /etc/ssl/key.pem &
```

### 添加文章
启动后，你只需将你的`*.md`文件放到`/usr/local/mknote/articles/`目录内即可，在这个目录中你可以创建多个目录作为文章的分类，需要注意的是你不能创建更深层级的目录，目前只支持一级子目录。
```
cd /usr/local/mknote/
mkdir articles/mknote
echo "# mknote" > articles/mknote/README.md
```

如果你的文章内引用了图片，只需将你的图放到`/usr/local/mknote/uploads/`目录下。
```
cp scenery.png /usr/local/mknote/uploads/
```

然后在文章内这样引用它。
```
![scenery](/uploads/scenery.png)
```

### 访问
现在可以在你的浏览器中打开`localhost`。

## 引用
* https://github.com/howeyc/fsnotify
* https://github.com/russross/blackfriday
* https://github.com/sindresorhus/github-markdown-css
