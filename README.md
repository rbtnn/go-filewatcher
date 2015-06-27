
# go-filewatcher

![](https://raw.githubusercontent.com/rbtnn/go-filewatcher/master/filewatcher.png)


## Configuration

```sh
go get github.com/mattn/go-colorable
go get github.com/shiena/ansicolor
make
```


## Command line options

Useful command line options:


### `-w {watchedDir}`

Watched directory (default: `.`).


### `-x {exclude}` _optional_

Execlude file name pattern (default: `.git,.hg,_svn`).


### `-d {depth}` _optional_

Maximum folder depth for which to watch (default: `0`).

