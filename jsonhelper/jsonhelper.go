package jsonhelper

import (
	"encoding/json"

	"github.com/lwh9346/MinecraftAutoUpdaterV2/filelist"
)

//JSONString 存有json格式数据的字符串
type JSONString = string

//GetFileListFromJSON 从json格式字符串中获取fileList
func GetFileListFromJSON(s JSONString) filelist.FileList {
	m := make(filelist.FileList)
	json.Unmarshal([]byte(s), &m)
	return m
}

//UpdateInfo 更新器更新所需要的信息
type UpdateInfo struct {
	GameVersion int      `json:"version"`
	IgnoreList  []string `json:"ignore_list"`
}

//LoadUpdateInfoFromJSON 从json格式字符串中获更新信息
func LoadUpdateInfoFromJSON(s JSONString) UpdateInfo {
	ui := UpdateInfo{}
	json.Unmarshal([]byte(s), &ui)
	return ui
}
