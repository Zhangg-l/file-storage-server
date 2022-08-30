package meta

import "go_code/project14/file-store-server/db"

type FileMeta struct {
	FileSha1 string
	FileName string
	FileSize int64
	UploadAt string
	Location string
}

var fileMetas map[string]FileMeta

func init() {
	fileMetas = make(map[string]FileMeta)
}

func UpdateFileMetas(f FileMeta) {
	fileMetas[f.FileSha1] = f
}

func UpLoadFileMetasOnDb(f FileMeta) bool {
	return db.OnUploadFileFinishied(f.FileSha1, f.FileName, f.Location, f.FileSize)
}

func GetFileMeta(fileSha string) FileMeta {
	return fileMetas[fileSha]
}

func GetFileMetaOnDb(fileSha string) (FileMeta, error) {
	tf, err := db.GetFileMeta(fileSha)
	if err != nil {
		return FileMeta{}, err
	}
	fm := FileMeta{
		FileSha1: tf.FileHash,
		FileName: tf.FileName.String,
		FileSize: tf.FileSize.Int64,
		Location: tf.FileAddr.String,
	}
	return fm ,nil
}

func RemoveFileMeta(fileSha string) {
	delete(fileMetas, fileSha)
}
