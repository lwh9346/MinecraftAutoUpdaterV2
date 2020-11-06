package main

import (
	"log"
	"os"
	"os/exec"
	"time"

	"github.com/lwh9346/MinecraftAutoUpdaterV2/download"

	"github.com/lwh9346/MinecraftAutoUpdaterV2/filelist"
	"github.com/lwh9346/MinecraftAutoUpdaterV2/jsonhelper"
)

const resourceURL = "https://minecraft-updater.oss-cn-shanghai.aliyuncs.com"

func main() {
	log.Println("MinecraftAutoUpdaterV2已启动")
	if len(os.Args) > 1 {
		switch os.Args[1] {
		case "init":
			initUpdateInfo()
		case "pack":
			makeUpdatePack()
		}
		return
	}
	autoUpdate()
	return
}

func initUpdateInfo() {
	log.Println("开始初始化更新包")
	os.RemoveAll("./pack")
	os.MkdirAll("./pack", os.ModePerm)
	filelist := filelist.GetFileList("./game")
	updateinfo := jsonhelper.UpdateInfo{GameVersion: 1}
	jsonhelper.WriteStringToFile("./pack/filelist.json", jsonhelper.GetJSONStringOfFileList(filelist))
	jsonhelper.WriteStringToFile("./pack/updateinfo.json", jsonhelper.GetJSONStringOfUpdateInfo(updateinfo))
	log.Println("更新包初始化完毕")
}

func makeUpdatePack() {
	log.Println("开始制作更新包")
	updateinfo := jsonhelper.LoadUpdateInfoFromJSON(jsonhelper.GetJSONStringByURL(resourceURL + "/updateinfo.json"))
	updateinfo.GameVersion++
	os.RemoveAll("./pack")
	os.MkdirAll("./pack", os.ModePerm)
	fl := filelist.GetFileList("./game")
	fl = filelist.IgnoreFileInFileList(updateinfo.IgnoreList, fl)
	jsonhelper.WriteStringToFile("./pack/filelist.json", jsonhelper.GetJSONStringOfFileList(fl))
	jsonhelper.WriteStringToFile("./pack/updateinfo.json", jsonhelper.GetJSONStringOfUpdateInfo(updateinfo))
	log.Println("更新包制作完毕")
}

func autoUpdate() {
	log.Println("开始自动更新")
	localUpdateInfo := jsonhelper.UpdateInfo{}
	_, e := os.Stat("./updateinfo.json")
	if e == nil {
		localUpdateInfo = jsonhelper.LoadUpdateInfoFromJSON(jsonhelper.GetJSONStringByFilePath("./updateinfo.json"))
	}
	updateinfo := jsonhelper.LoadUpdateInfoFromJSON(jsonhelper.GetJSONStringByURL(resourceURL + "/updateinfo.json"))
	log.Printf("已获取最新版本信息，当前版本:%d，最新版本:%d\n", localUpdateInfo.GameVersion, updateinfo.GameVersion)
	if updateinfo.GameVersion > localUpdateInfo.GameVersion {
		log.Println("开始更新所需下载文件并删除旧文件")
		os.MkdirAll("./game", os.ModePerm)
		oldFileList := filelist.GetFileList("./game")
		newFileList := jsonhelper.GetFileListFromJSON(jsonhelper.GetJSONStringByURL(resourceURL + "/filelist.json"))
		oldFileList = filelist.IgnoreFileInFileList(updateinfo.IgnoreList, oldFileList)
		newFileList = filelist.IgnoreFileInFileList(updateinfo.IgnoreList, newFileList)
		surp, lack := filelist.CompareFileList(oldFileList, newFileList)
		for path, _ := range surp {
			os.Remove(path)
		}
		failed := download.DownloadAndCheckFilesInFileList(resourceURL, lack)
		log.Println("下载完毕")
		if failed == 0 {
			jsonhelper.WriteStringToFile("./updateinfo.json", jsonhelper.GetJSONStringOfUpdateInfo(updateinfo))
		}
	} else {
		log.Println("已是最新版")
	}
	launchGameLauncher()
}

func launchGameLauncher() {
	log.Println("正在启动游戏启动器，请不要关闭更新器窗口")
	cmd := exec.Command("java", "-jar", "Launcher.jar")
	cmd.Dir = "./game"
	//cmd.Stdout = os.Stdout
	err := cmd.Start()
	if err != nil {
		log.Println("启动失败，你可能没安装java")
		time.Sleep(60 * time.Second)
	}
}