package download

import (
	"io"
	"log"
	"net/http"
	"os"
	fp "path/filepath"
	"strings"

	"github.com/lwh9346/MinecraftAutoUpdaterV2/filelist"
)

func downloadFile(url, destDir string) error {
	destFile, err := os.Create(destDir)
	defer destFile.Close()
	if err != nil {
		return err
	}
	var res *http.Response
	res, err = http.Get(url)
	if err != nil {
		return err
	}
	_, err = io.Copy(destFile, res.Body)
	return err
}

func downloadFileAndCheck(url, destDir, hash string, limitor, callback chan (int)) error {
	limitor <- 1
	_, e := os.Stat(destDir)
	if e == nil && hash == filelist.GetHash(destDir) {
		callback <- 0
		return nil
	} else {
		os.Remove(destDir)
	}
	for i := 0; ; i++ {
		err := downloadFile(url, destDir)
		if err == nil && hash == filelist.GetHash(destDir) {
			callback <- 0
			return nil
		}
		if i > 10 {
			callback <- 1
			return err
		}
	}
}

func DownloadAndCheckFilesInFileList(rootURL string, filelist filelist.FileList) {
	nFiles := len(filelist)
	var succeed, failed int
	limitor := make(chan (int), 16)
	callback := make(chan (int))
	for filepath, filehash := range filelist {
		go downloadFileAndCheck(strings.Join([]string{rootURL, filepath}, "/"), fp.FromSlash(filepath), filehash, limitor, callback)
	}
	if nFiles == 0 {
		return
	}
	for signal := range callback {
		<-limitor
		if signal == 0 {
			succeed++
		} else {
			failed++
		}
		uncompleted := nFiles - succeed - failed
		if uncompleted == 0 {
			close(callback)
		}
		if uncompleted%20 == 0 {
			log.Printf("下载成功:%v,下载失败:%v,尚未下载:%v\n", succeed, failed, uncompleted)
		}
	}

}
