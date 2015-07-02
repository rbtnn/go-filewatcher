package main

import (
	"crypto/sha256"
	"encoding/hex"
	"flag"
	"fmt"
	"github.com/mattn/go-colorable"
	"github.com/shiena/ansicolor"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

type EventType string

const (
	Added   EventType = "Added"
	Updated           = "Updated"
	Deleted           = "Deleted"
)

type TT struct {
	Path             string
	LastModifiedTime time.Time
	Size             int64
	Hash             string
}

type Table map[string]TT

type Arguments struct {
	table      Table
	watchedDir string
	exclude    []string
	depth      int
}

func main() {
	watchedDir_flag := flag.String("w", ".", "")
	exclude_flag := flag.String("x", ".git,.hg,_svn", "")
	depth_flag := flag.Int("d", 0, "")

	flag.Parse()

	args := Arguments{
		make(Table, 0),
		*watchedDir_flag,
		strings.Split(*exclude_flag, ","),
		*depth_flag,
	}

	_, err := os.Stat(args.watchedDir)
	if err == nil {
		for {
			time.Sleep(1000 * time.Millisecond)
			listFiles(&args, args.watchedDir)
			for k, v := range args.table {
				_, err := os.Stat(v.Path)
				if err != nil {
					v.Message(Deleted, &TT{})
					delete(args.table, k)
				}
			}
		}
	}
}

func (t *TT) Message(e EventType, prev *TT) {
	w := ansicolor.NewAnsiColorWriter(colorable.NewColorableStdout())
	text := "[" + time.Now().Format(time.Stamp) + "] " + string(e) + ": "
	switch e {
	case Added:
		// Yellow
		size := strconv.FormatInt(t.Size, 10) + " Bytes"
		text += t.Path + " (" + size + ", " + t.Hash + ")"
		fmt.Fprintf(w, "%s%s%s\n", "\x1b[33m", text, "\x1b[0m")
	case Updated:
		// Green
		size := strconv.FormatInt(prev.Size, 10) + " -> " + strconv.FormatInt(t.Size, 10) + " Bytes"
		text += t.Path + " (" + size + ", " + prev.Hash + " -> " + t.Hash + ")"
		fmt.Fprintf(w, "%s%s%s\n", "\x1b[32m", text, "\x1b[0m")
	case Deleted:
		// Red
		size := strconv.FormatInt(t.Size, 10) + " Bytes"
		text += t.Path + " (" + size + ", " + t.Hash + ")"
		fmt.Fprintf(w, "%s%s%s\n", "\x1b[31m", text, "\x1b[0m")
	}
}

func getSha256(path string) string {
	hasher := sha256.New()
	f, err := os.Open(path)
	if err != nil {
		return ""
	}
	defer f.Close()
	if _, err := io.Copy(hasher, f); err != nil {
		return ""
	}
	return hex.EncodeToString(hasher.Sum(nil))[:16]
}

func listFiles(args *Arguments, path string) {
	fis, err := ioutil.ReadDir(path)
	if err == nil {
		if 0 <= args.depth {
			args.depth--
			for _, fi := range fis {
				ok := true
				for _, pattern := range args.exclude {
					if matched, _ := filepath.Match(pattern, fi.Name()); matched {
						ok = false
						break
					}
				}
				if ok {
					if fi.IsDir() {
						listFiles(args, filepath.Join(path, fi.Name()))
					} else {
						check(args, filepath.Join(path, fi.Name()))
					}
				}
			}
			args.depth++
		}
	}
}

func check(args *Arguments, path string) {
	file, err := os.Stat(path)
	if err == nil {
		lastModifiedTime := file.ModTime()
		prev, ok := args.table[path]
		if ok {
			if prev.LastModifiedTime != lastModifiedTime {
				hash := getSha256(path)
				if prev.Hash != hash {
					v := TT{path, lastModifiedTime, file.Size(), hash}
					v.Message(Updated, &prev)
					args.table[path] = v
				}
			}
		} else {
			hash := getSha256(path)
			v := TT{path, lastModifiedTime, file.Size(), hash}
			v.Message(Added, &prev)
			args.table[path] = v
		}
	}
}
