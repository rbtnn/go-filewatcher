
# go-filewatcher

![](https://raw.githubusercontent.com/rbtnn/go-filewatcher/master/filewatcher.png)


## Configuration

```sh
go get github.com/mattn/go-colorable
go get github.com/shiena/ansicolor
make
```


## Usage

__Windows__

```
filewatcher.exe -r {rootdir} [-i {ignorelist}] [-d {depth}]
```

__Others__

```
filewatcher -r {rootdir} [-i {ignorelist}] [-d {depth}]
```


## Option


### `-r {rootdir}`

default: `-r .`


### `-i {ignorelist}` optional

default: `-i ".git,.hg,_svn"`


### `-d {depth}` optional

default: `-d 0`

