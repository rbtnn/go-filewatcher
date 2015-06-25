package main

import (
	"flag"
	"fmt"
	"github.com/mattn/go-colorable"
	"github.com/shiena/ansicolor"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
	"time"
)

type EventType string

const (
	Added   EventType = "Added"
	Updated           = "Updated"
	Deleted           = "Deleted"
)

type TT struct {
	AbstractPath     string
	RelativePath     string
	LastModifiedTime time.Time
	Size             int64
}

func (t *TT) Message(e EventType) {
	w := ansicolor.NewAnsiColorWriter(colorable.NewColorableStdout())
	size := strconv.FormatInt(t.Size, 10) + " Bytes"
	text := "[" + t.LastModifiedTime.Format(time.Stamp) + "] " + string(e) + ": " + t.RelativePath + " (" + size + ")"
	switch e {
	case Added:
		// Yellow
		fmt.Fprintf(w, "%s%s%s\n", "\x1b[33m", text, "\x1b[0m")
	case Updated:
		// Green
		fmt.Fprintf(w, "%s%s%s\n", "\x1b[32m", text, "\x1b[0m")
	case Deleted:
		// Red
		fmt.Fprintf(w, "%s%s%s\n", "\x1b[31m", text, "\x1b[0m")
	}
}

type Table map[string]TT

func main() {
	root_flag := flag.String("d", "", "rootdir")
	flag.Parse()
	root := *root_flag
	_, err := os.Stat(root)
	if err == nil {
		table := make(Table, 0)
		for {
			time.Sleep(1000 * time.Millisecond)
			listFiles(root, root, table)
			for k, v := range table {
				_, err := os.Stat(v.AbstractPath)
				if err != nil {
					v.Message(Deleted)
					delete(table, k)
				}
			}
		}
	}
}

func listFiles(rootPath string, searchPath string, table Table) {
	fis, err := ioutil.ReadDir(searchPath)
	if err == nil {
		for _, fi := range fis {
			fullPath := filepath.Join(searchPath, fi.Name())
			if fi.IsDir() {
				listFiles(rootPath, fullPath, table)
			} else {
				check(fullPath, fi.Name(), table)
			}
		}
	}
}

func check(abstractPath string, relativePath string, table Table) {
	file, err := os.Stat(abstractPath)
	if err == nil {
		lastModifiedTime := file.ModTime()
		prev, ok := table[abstractPath]
		if ok {
			if prev.LastModifiedTime != lastModifiedTime {
				v := TT{abstractPath, relativePath, lastModifiedTime, file.Size()}
				v.Message(Updated)
				table[abstractPath] = v
			}
		} else {
			v := TT{abstractPath, relativePath, lastModifiedTime, file.Size()}
			v.Message(Added)
			table[abstractPath] = v
		}
	}
}
