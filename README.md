# Marknote
mknote is a simple and quick blogging system, No database, No login required, Only standerd documents are needed.

## Documents
* [中文文档](https://github.com/sycki/mknote/blob/master/README_CH.md)
* [English doc](https://github.com/sycki/mknote)

## Quick start
Before you start, please make sure you have installed golang and configured the GOPATH environment variable.

### Install mknote
Where /usr/local/mknote refers to the installation path and must be a directory that does not exist.
```
go get github.com/sycki/mknote
$GOPATH/src/github.com/sycki/mknote/build.sh install /usr/local/mknote
```

Confirm installed.
```
cd /usr/local/mknote/ && ls
articles  bin  conf  f  static
```

### Launch mknote
Usually started in https mode, you need to specify your certificate file and the final domain name.
```
bin/mknote \
--hostname blog.domain.com \
--tls=true \
--tls-cert /etc/ssl/cert.pem \
--tls-key /etc/ssl/key.pem
```

Or start it in http mode:
```
bin/mknote
```

mknote provides many useful options, view all options with the following command:
```
bin/mknote --help
```

### Writing an article
Create a new article category and write your first article.
```
mkdir articles/java
echo '# First article' > articles/java/first.md
```

Now you can access `http://localhost` in your browser.

## File download server
If you need to reference the image in your article, just put your image in the /usr/local/mknote/f/ directory, which can be specified when starting mknote.
```
cp scenery.png /usr/local/mknote/f/
```

Then reference it in the article
```
![scenery](/f/scenery.png)
```

In fact, you can place any file in this directory, so that others can download the file anywhere. This function is very useful for many people. You can add as many directories as you like in the /usr/local/mknote/f/ directory. To distinguish your files.

## Debug
Open debug feature
```
curl -X POST -H "<your_header_key>: <value>" https://sycki.com/v1/manage/pprof/open
```

Analyse your mknote using the go tool `go profile`
```
go tool pprof http://sycki.com:8000/debug/pprof/profile
```

Close debug feature
```
curl -X POST -H "<your_header_key>: <value>" https://sycki.com/v1/manage/pprof/close
```

## Reference
* https://github.com/howeyc/fsnotify
* https://github.com/russross/blackfriday
* https://github.com/sindresorhus/github-markdown-css

