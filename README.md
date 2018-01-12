# Marknote
mknote is a simple and quick blogging system, No database, No login required, Only standerd documents are needed.

## Documents
* [中文文档](https://github.com/sycki/mknote/blob/master/README-CH.md)
* [English doc](https://github.com/sycki/mknote)

## Development guide
### Clone source code
You can fork the project first, Also can direct clone it.
```
cd $GOPATH/src
git clone https://github.com/sycki/mknote.git
```

### Development
Import the project to your IDE, Edition it.

### Build and run
```
cd $GOPATH/src/mknote
go build -v mknote
sudo ./mknote
```

## Usage guide
### Download
Download latest binary tarball at release page.
[https://github.com/sycki/mknote/releases](https://github.com/sycki/mknote/releases)

### Uncompress
```
mkdir /usr/local/mknote/
tar -xf mknote-v2.2.tar -C /usr/local/mknote/
```

### Run
```
/usr/local/sycki-mknote/mknote &
```

### Run with TLS
Specify your cert and key file by options `--tls-cert` and `--tls-key`, If you want start it with TLS mode.
```
/usr/local/sycki-mknote/mknote \
--tls-cert /etc/ssl/cert.pem \
--tls-key /etc/ssl/key.pem &
```

### Add articles
Copy your *.md files to `/usr/local/mknote/articles/`, You can create subdirectory in `articles`, And at most 1 level.
```
cd /usr/local/mknote/
mkdir articles/mknote
echo "# mknote" > articles/mknote/README.md
```

Copy your images to `/usr/local/mknote/uploads/`.
```
cp scenery.png /usr/local/mknote/uploads/
```

After then refrence they in your articles.
```
![scenery](/uploads/scenery.png)
```

### Visit
Visit `localhost` in your browser.

## Refrence
* https://github.com/howeyc/fsnotify
* https://github.com/russross/blackfriday
* https://github.com/sindresorhus/github-markdown-css
