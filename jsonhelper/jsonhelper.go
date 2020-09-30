package jsonhelper

import (
	"encoding/json"

	"github.com/lwh9346/MinecraftAutoUpdaterV2/filelist"
)

type JSONString = string

func GetFileListFromJSON(s JSONString) filelist.FileList {
	m := make(filelist.FileList)
	json.Unmarshal([]byte(s), &m)
	return m
}

type UpdateInfo struct {
	GameVersion int      `json:"version"`
	IgnoreList  []string `json:"ignore_list"`
}

func LoadUpdateInfoFromJSON(s JSONString) UpdateInfo {
	ui := UpdateInfo{}
	json.Unmarshal([]byte(s), &ui)
	return ui
}
