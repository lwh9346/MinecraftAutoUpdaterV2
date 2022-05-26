package main

import (
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"time"

	"github.com/gonutz/w32/v2"

	"mau2/filelist"
	"mau2/updateinfo"
	"mau2/utils"
)

const resourceURL = "https://minecraft-updater.oss-accelerate.aliyuncs.com"

func main() {
	log.Println("MinecraftAutoUpdaterV2已启动")
	if len(os.Args) > 1 {
		switch os.Args[1] {
		case "init":
			initUpdateInfo()
		case "pack":
			makeUpdatePack()
		case "finishselfupdate":
			finishSelfUpdate()
		case "hash":
			selfHash := filelist.GetHash(os.Args[0])
			utils.WriteStringToFile("./updaterhash", selfHash)
			log.Println(selfHash)
		case "launch":
			launchGameLauncher()
		}
		return
	}
	autoUpdate()
}

func initUpdateInfo() {
	log.Println("开始初始化更新包")
	os.RemoveAll("./pack")
	os.MkdirAll("./pack", os.ModePerm)
	fl := filelist.GetFileList("./game")
	ui := updateinfo.UpdateInfo{GameVersion: 1}
	utils.WriteStringToFile("./pack/filelist.json", updateinfo.ToJSON(ui))
	utils.WriteStringToFile("./pack/updateinfo.json", filelist.ToJSON(fl))
	log.Println("更新包初始化完毕")
}

func makeUpdatePack() {
	log.Println("开始制作更新包")
	ui := updateinfo.FromJSON(utils.ReadStringFromURL(resourceURL + "/updateinfo.json"))
	ui.GameVersion++
	os.RemoveAll("./pack")
	os.MkdirAll("./pack", os.ModePerm)
	fl := filelist.GetFileList("./game")
	filelist.Ignore(ui.IgnoreList, fl)
	utils.WriteStringToFile("./pack/filelist.json", filelist.ToJSON(fl))
	utils.WriteStringToFile("./pack/updateinfo.json", updateinfo.ToJSON((ui)))
	log.Println("更新包制作完毕")
}

func autoUpdate() {
	selfUpdate()
	log.Println("开始自动更新")
	localUpdateInfo := updateinfo.UpdateInfo{}
	_, e := os.Stat("./updateinfo.json")
	if e == nil {
		localUpdateInfo = updateinfo.FromJSON(utils.ReadStringFromFile("./updateinfo.json"))
	}
	ui := updateinfo.FromJSON(utils.ReadStringFromURL(resourceURL + "/updateinfo.json"))
	log.Printf("已获取最新版本信息，当前版本:%d，最新版本:%d\n", localUpdateInfo.GameVersion, ui.GameVersion)
	if ui.GameVersion > localUpdateInfo.GameVersion {
		log.Println("开始更新所需下载文件并删除旧文件")
		os.MkdirAll("./game", os.ModePerm)
		oldFileList := filelist.GetFileList("./game")
		newFileList := filelist.FromJSON(utils.ReadStringFromURL(resourceURL + "/filelist.json"))
		filelist.Ignore(ui.IgnoreList, oldFileList)
		filelist.Ignore(ui.IgnoreList, newFileList)
		surp, lack := filelist.CompareFileList(oldFileList, newFileList)
		for path := range surp {
			os.Remove(path)
		}
		utils.RemoveEmptyDirectories("./game")
		failed := utils.DownloadAndCheckFilesInFileList(resourceURL, lack)
		log.Println("下载完毕")
		if failed == 0 {
			utils.WriteStringToFile("./updateinfo.json", updateinfo.ToJSON(ui))
		}
	} else {
		log.Println("已是最新版")
	}
	launchGameLauncher()
}

func launchGameLauncher() {
	log.Println("正在启动游戏启动器，请不要关闭更新器窗口")
	launcherType := "none"
	if _, e := os.Stat("./game/Launcher.jar"); e == nil {
		launcherType = "jar"
	}
	if _, e := os.Stat("./game/Launcher.exe"); e == nil {
		launcherType = "exe"
	}
	var cmd *exec.Cmd
	switch launcherType {
	case "jar":
		cmd = exec.Command("./jre8/bin/java.exe", "-jar", "Launcher.jar")
	case "exe":
		cmd = exec.Command("./Launcher.exe")
	case "none":
		log.Println("启动失败，找不到启动器，请联系管理员")
		time.Sleep(60 * time.Second)
		os.Exit(0)
	}
	cmd.Dir = "./game"
	err := cmd.Start()
	if err != nil {
		log.Println("启动失败，请联系管理员")
		time.Sleep(60 * time.Second)
	}
	//使用win32api关闭命令行窗口
	console := w32.GetConsoleWindow()
	w32.ShowWindow(console, w32.SW_HIDE)
}

func selfUpdate() {
	log.Println("开始更新器更新检查")
	selfHash := filelist.GetHash(os.Args[0])
	latestHash := utils.ReadStringFromURL(resourceURL + "/updaterhash")
	if latestHash == selfHash {
		log.Println("更新器已是最新版")
		return
	}
	log.Printf("当前程序hash：%s\n", selfHash)
	log.Printf("最新程序hash：%s\n", latestHash)
	log.Println("开始自我更新")
	newProgramURL := resourceURL + "/program/AutoUpdater.exe"
	exeName := filepath.Join(filepath.Dir(os.Args[0]), "SelfUpdater.exe")
	os.Remove(exeName)
	utils.DownloadFile(resourceURL+"/program/SelfUpdater.exe", exeName)
	cmd := exec.Command(exeName, os.Args[0], newProgramURL)
	cmd.Dir = filepath.Dir(os.Args[0])
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Start()
	os.Exit(0)
}

func finishSelfUpdate() {
	os.Remove(filepath.Join(filepath.Dir(os.Args[0]), "SelfUpdater.exe"))
	log.Println("自我更新完成")
	autoUpdate()
}
