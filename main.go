package main

import (
	"log"

	"github.com/lwh9346/MinecraftAutoUpdaterV2/filelist"
)

func main() {
	log.Println("MinecraftAutoUpdaterV2已启动...")
	filelist := filelist.GetFileList("./game")
	log.Printf("%v\n", filelist)
}
