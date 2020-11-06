package jsonhelper

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/lwh9346/MinecraftAutoUpdaterV2/filelist"
)

//GetJSONStringByURL 从指定地址以get形式获取string
func GetJSONStringByURL(url string) JSONString {
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

//GetJSONStringByFilePath 从指定本地地址以get形式获取string
func GetJSONStringByFilePath(filepath string) JSONString {
	b, err := ioutil.ReadFile(filepath)
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
	defer f.Close()
	if e != nil {
		return e
	}
	io.WriteString(f, s)
	return nil
}

//GetJSONStringOfUpdateInfo 顾名思义
func GetJSONStringOfUpdateInfo(ui UpdateInfo) JSONString {
	d, _ := json.Marshal(ui)
	return string(d)
}

//GetJSONStringOfFileList 顾名思义
func GetJSONStringOfFileList(fl filelist.FileList) JSONString {
	d, _ := json.Marshal(fl)
	return string(d)
}
