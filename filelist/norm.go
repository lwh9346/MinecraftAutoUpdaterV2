package filelist

import "path/filepath"

//NormcaseFilelist 将目录名称中的slash转换为当前系统分隔符，注意该方法仅用于在下载文件之前对filelist进行操作
func NormcaseFilelist(fileList FileList) FileList {
	normFileList := make(FileList)
	for k, v := range fileList {
		normFileList[filepath.FromSlash(k)] = v
	}
	return normFileList
}

//toSlashFilelist 将目录名称中的分隔符转换为slash
func toSlashFilelist(fileList FileList) FileList {
	slashFileList := make(FileList)
	for k, v := range fileList {
		slashFileList[filepath.ToSlash(k)] = v
	}
	return slashFileList
}
