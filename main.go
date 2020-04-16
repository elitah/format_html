package main

import (
	"bytes"
	"crypto/md5"
	"flag"
	"fmt"
	"html/template"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"time"
)

func StringSortAdd(list *[]string, str string) {
	if "" != str {
		for _, s := range *list {
			if str == s {
				return
			}
		}
		*list = append(*list, str)
	}
}

func main() {
	var help bool
	var write bool

	var indexFile string
	var jsDir string
	var jsList string

	flag.BoolVar(&help, "h", false, "This Help.")
	flag.StringVar(&indexFile, "i", "", "your index.html path")
	flag.StringVar(&jsDir, "j", "", "your js path")
	flag.StringVar(&jsList, "l", "", "you js list")
	flag.BoolVar(&write, "w", false, "enable write")

	flag.Parse()

	if help {
		// for help
	} else if "" != indexFile && "" != jsDir && "" != jsList {
		var list []string

		if _path, err := filepath.Abs(indexFile); nil == err {
			indexFile = _path
		}

		if _path, err := filepath.Abs(jsDir); nil == err {
			jsDir = _path
		}

		//fmt.Println("!!!", indexFile, jsDir, jsList)

		// 遍历js
		if _list, err := filepath.Glob(jsDir + "/*.js"); nil == err {
			for _, item := range strings.Split(jsList, ";") {
				if 3 < len(item) {
					StringSortAdd(&list, item[:len(item)-3])
				}
			}
			for _, item := range _list {
				item = filepath.Base(item)
				if 3 < len(item) {
					StringSortAdd(&list, item[:len(item)-3])
				}
			}
		} else {
			fmt.Println(err)
		}

		if 0 < len(list) {
			if data, err := ioutil.ReadFile(indexFile); nil == err {
				if t := template.New("index.tpl"); nil != t {
					// 加载模板
					if _, err = t.Parse(string(data)); nil == err {
						var b bytes.Buffer
						var buffer [1024]byte
						var path string
						// md5
						h := md5.New()
						// 计算js校验和同时修改文件名
						for i, item := range list {
							path = fmt.Sprintf("%s/%s.js", jsDir, item)
							if file, err := os.Open(path); nil == err {
								h.Reset()
								io.CopyBuffer(h, file, buffer[:])
								list[i] = fmt.Sprintf("%s.%x.js", item, h.Sum(nil)[:])
								if write {
									os.Rename(path, fmt.Sprintf("%s/%s", jsDir, list[i]))
								}
							}
						}
						// 执行模板
						if err = t.Execute(&b, struct {
							List        []string
							CurrentYear int
						}{
							List:        list,
							CurrentYear: time.Now().UTC().Add(8 * time.Hour).Year(),
						}); nil == err {
							if write {
								ioutil.WriteFile(indexFile, b.Bytes(), 0644)
							} else {
								fmt.Println(b.String())
							}
							return
						} else {
							fmt.Println(err)
						}
					} else {
						fmt.Println(err)
					}
				} else {
					fmt.Println("template.New() return failed")
				}
			} else {
				fmt.Println(err)
			}
		}

		os.Exit(-1)

		return
	}

	flag.Usage()
}
