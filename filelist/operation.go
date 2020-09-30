package filelist

func CompareFileList(oldFileList, newFileList FileList) (surp, lack FileList) {
	su := make(FileList)
	la := make(FileList)
	for k, vo := range oldFileList {
		vn, exists := newFileList[k]
		if (!exists) || vn != vo { //newFileList里不存在或者hash对不上的文件被认为是多余的
			su[k] = vo
		}
	}
	for k, vn := range newFileList {
		vo, exists := oldFileList[k]
		if (!exists) || vn != vo { //oldFileList里不存在或者hash对不上的文件被认为是缺失的
			la[k] = vn
		}
	}
	return su, la
}

func IgnoreFileInFileList(ignoreList []string, fileList FileList) FileList {
	del := make([]string, len(ignoreList))
	for kf := range fileList {
		for _, ki := range ignoreList {
			kfn := []rune(kf)
			kin := []rune(ki)
			if len(kin) <= len(kfn) {
				if ki == string(kfn[:len(kin)]) {
					del = append(del, kf)
				}
			}
		}
	}
	for _, d := range del {
		if d != "" {
			delete(fileList, d)
		}
	}
	return fileList
}
