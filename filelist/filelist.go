package filelist

import (
	"crypto/sha1"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"
)

type FileList = map[string]string //在传输、保存的过程中都使用slash作为分隔符

func GetHash(fp string) string {
	sha := sha1.New()
	f, _ := os.Open(fp)
	defer f.Close()
	io.Copy(sha, f)
	return strings.ToLower(fmt.Sprintf("%X", sha.Sum(nil)))
}

type FileListElement struct {
	FilePath string
	FileHash string
}

func GetFileList(path string) FileList {
	nFiles := 0
	fileListElementReceiver := make(chan FileListElement)
	filepath.Walk(path, func(p string, fi os.FileInfo, e error) error {
		if e != nil {
			log.Printf("%v\n", e)
			return e
		}
		if !fi.IsDir() {
			nFiles += 1
			p = filepath.ToSlash(p)
			go func() {
				fileListElementReceiver <- FileListElement{FileHash: GetHash(p), FilePath: p}
			}()
		}
		return nil
	})
	fileList := make(FileList)
	if nFiles == 0 {
		close(fileListElementReceiver)
	}
	for fileListElement := range fileListElementReceiver {
		fileList[fileListElement.FilePath] = fileListElement.FileHash
		nFiles--
		if nFiles == 0 {
			close(fileListElementReceiver)
		}
	}
	return fileList
}

func CompareFileList(old, new FileList) (surp, lack FileList) {
	su := make(FileList)
	la := make(FileList)
	for k, vo := range old {
		vn, exists := new[k]
		if (!exists) || vn != vo { //new里不存在或者hash对不上的文件被认为是多余的
			su[k] = vo
		}
	}
	for k, vn := range new {
		vo, exists := old[k]
		if (!exists) || vn != vo { //old里不存在或者hash对不上的文件被认为是缺失的
			la[k] = vn
		}
	}
	return su, la
}

func Ignore(ignoreList []string, fileList FileList) {
	del := make([]string, len(ignoreList))
	for kf := range fileList {
		for _, ki := range ignoreList {
			if strings.HasPrefix(kf, ki) {
				del = append(del, kf)
			}
		}
	}
	for _, d := range del {
		if d != "" {
			delete(fileList, d)
		}
	}
}

func FromJSON(s string) FileList {
	m := make(FileList)
	json.Unmarshal([]byte(s), &m)
	return m
}

func ToJSON(fl FileList) string {
	d, _ := json.Marshal(fl)
	return string(d)
}
