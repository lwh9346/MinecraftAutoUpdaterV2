package utils

import (
	"io"
	"io/ioutil"
	"net/http"
	"os"
)

//ReadStringFromURL 从指定地址以get形式获取string
func ReadStringFromURL(url string) string {
	r, e := http.Get(url)
	if e != nil {
		return ""
	}
	b, e := ioutil.ReadAll(r.Body)
	if e != nil {
		return ""
	}
	return string(b)
}

//ReadStringFromFile 从指定本地文件获取string
func ReadStringFromFile(file string) string {
	b, err := ioutil.ReadFile(file)
	if err != nil {
		return ""
	}
	return string(b)
}

//WriteStringToFile 顾名思义
func WriteStringToFile(file, s string) error {
	var e error
	_, e = os.Stat(file)
	if e == nil {
		e = os.Remove(file)
		if e != nil {
			return e
		}
	}
	f, e := os.Create(file)
	if e != nil {
		return e
	}
	defer f.Close()
	io.WriteString(f, s)
	return nil
}
