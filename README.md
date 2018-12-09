# Marknote
Mknote is an easy-to-use blogging system with no database and no need to log in. Just specify the directory where your price reduction files are located, mknote will automatically discover them and generate pages.

## Documents
* [中文文档](https://github.com/sycki/mknote/blob/master/README_ZH.md)
* [English doc](https://github.com/sycki/mknote)

## Quick start
### Install mknote
Go to the [Publish Page](https://github.com/sycki/mknote/releases) to download the latest version of the package and extract it to your favorite directory.
```
tar -zxf mknote-<version>.tar.gz
cd mknote-<version>
```

### Launch mknote
Usually started in https mode, you need to specify your certificate file and the final domain name, which is used to redirect http requests to https:
```
bin/mknote \
--tls=true \
--tls-cert /etc/ssl/cert.pem \
--tls-key /etc/ssl/key.pem \
--hostname blog.domain.com
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
curl -X POST -H "<your_header_key>: <value>" https://<hostname>/v1/manage/pprof/open
```

Analyse your mknote using the go tool `go profile`
```
go tool pprof http://<hostname>:8000/debug/pprof/profile
```

Close debug feature
```
curl -X POST -H "<your_header_key>: <value>" https://<hostname>/v1/manage/pprof/close
```
