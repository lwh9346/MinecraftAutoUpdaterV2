package updateinfo

import (
	"encoding/json"
	"errors"
	"mau2/config"
	"mau2/utils"
	"os"
)

type Fix struct {
	Operation string `json:"operation"` //remove change(add)
	Target    string `json:"target"`    //要修改的文件
	Source    string `json:"source"`    //op为remove时留空
}

// UpdateInfo 更新器更新所需要的信息
type UpdateInfo struct {
	GameVersion int      `json:"version"`
	IgnoreList  []string `json:"ignore_list"`
	FixVersion  int      `json:"fix_version"`
	Fixs        []Fix    `json:"fixs"`
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

func (hf *Fix) Do() error {
	switch hf.Operation {
	case "remove":
		os.Remove(hf.Target)
		return nil
	case "change":
		os.Remove(hf.Target)
		return utils.DownloadFile(config.ResourceURL+"/fix/"+hf.Source, hf.Target)
	default:
		return errors.New("invalid operation")
	}
}
