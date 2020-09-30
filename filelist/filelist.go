package filelist

type FilePath = string
type FileHash = string
type FileList = map[FilePath]FileHash //在传输、保存的过程中都使用slash作为分隔符
type FileListElement struct {
	FilePath FilePath
	FileHash FileHash
}
