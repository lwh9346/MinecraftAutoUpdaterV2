package updateinfo

import "encoding/json"

//UpdateInfo 更新器更新所需要的信息
type UpdateInfo struct {
	GameVersion int      `json:"version"`
	IgnoreList  []string `json:"ignore_list"`
}

func FromJSON(s string) UpdateInfo {
	var m UpdateInfo
	json.Unmarshal([]byte(s), &m)
	return m
}

func ToJSON(ui UpdateInfo) string {
	d, _ := json.Marshal(ui)
	return string(d)
}
