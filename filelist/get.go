package filelist

import (
	"crypto/sha1"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"
)

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
	return toSlashFilelist(fileList)
}

func GetHash(filepath FilePath) FileHash {
	sha := sha1.New()
	f, _ := os.Open(filepath)
	defer f.Close()
	io.Copy(sha, f)
	return strings.ToLower(fmt.Sprintf("%X", sha.Sum(nil)))
}
