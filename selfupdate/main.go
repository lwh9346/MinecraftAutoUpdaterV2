package main

import (
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"time"

	"github.com/lwh9346/MinecraftAutoUpdaterV2/download"
)

func main() {
	fileToUpdate := os.Args[1]
	log.Printf("正在删除：%s\n", fileToUpdate)
	err := os.Remove(fileToUpdate)
	if err != nil {
		log.Fatalf("自我更新失败，删除文件时出现错误：%s\n", err.Error())
	}
	log.Println("成功删除文件")
	fileToDownload := os.Args[2]
	log.Println("开始下载新版更新器")
	err = download.DownloadFile(fileToDownload, fileToUpdate)
	if err != nil {
		log.Printf("自我更新失败，无法下载文件：%s\n", fileToDownload)
		log.Println("你需要重新下载更新器了")
		time.Sleep(time.Second * 30)
		os.Exit(1)
	}
	log.Println("下载完毕")
	exeName := fileToUpdate
	cmd := exec.Command(exeName, "finishselfupdate")
	cmd.Dir = filepath.Dir(os.Args[0])
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	log.Println("正在启动新版更新器")
	cmd.Start()
}
