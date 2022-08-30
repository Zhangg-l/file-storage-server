package handler

import (
	"encoding/json"
	"fmt"
	"go_code/project14/file-store-server/meta"
	"go_code/project14/file-store-server/util"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
)

func UploadHandler(c *gin.Context) {
	data, err := ioutil.ReadFile("./static/view/upload.html")
	if err != nil {
		c.String(404, `网页不存在`)
		return
	}
	c.Data(http.StatusOK, "text/html; charset=utf-8", data)
}

func DoUploadHandler(c *gin.Context) {
	var fileMeta meta.FileMeta
	var newFile *os.File
	var msg string
	var err error

	// parse upload 　form file
	file, fHead, err := c.Request.FormFile("file")
	defer file.Close()
	if err != nil {
		msg = fmt.Sprintf("failed to FormFile file:%s", err)
		goto ERR

	}
	// create new file to save user data
	newFile, err = os.Create("/tmp/" + fHead.Filename)
	defer func() {
		if newFile != nil {
			newFile.Close()
		}
	}()
	if err != nil {
		msg = fmt.Sprintf("failed to Create file:%s", err)
		goto ERR
	}

	fileMeta = meta.FileMeta{
		UploadAt: time.Now().Format("2006-01-02 15:04:05"),
		FileName: fHead.Filename,
		Location: "/tmp/" + fHead.Filename,
	}
	fileMeta.FileSize, err = io.Copy(newFile, file)
	if err != nil {
		msg = fmt.Sprintln("failed to Copy new file & file:", err)
		goto ERR
	}
	//
	newFile.Seek(0, 0)
	fileMeta.FileSha1 = util.FileSha1(newFile)
	// redirect to success file
	if !meta.UpLoadFileMetasOnDb(fileMeta) {
		msg = fmt.Sprintf("failed to UpLoadFileMetasOnDb")
		goto ERR
	}

	c.Redirect(http.StatusFound, "/file/upload/suc")
	return
ERR:
	c.JSON(http.StatusOK, gin.H{
		"code": "-1",
		"msg":  msg,
	})
	return
}

func UploadSucHandler(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"code": "0",
		"msg":  "upload file success!",
	})
}

func GetFileMetaHandler(c *gin.Context) {
	var err error
	var msg string
	var fmeta meta.FileMeta
	var data []byte
	filesha1 := c.Request.FormValue("filehash")
	fmeta, err = meta.GetFileMetaOnDb(filesha1)
	if err != nil {
		msg = fmt.Sprintln("meta.GetFileMetaOnDb:", err)
		goto ERR
	}
	data, err = json.Marshal(fmeta)
	if err != nil {
		msg = fmt.Sprintln("failed to Marshal:", err)
		goto ERR
	}

	c.JSON(http.StatusOK, gin.H{
		"code": "0",
		"msg":  data,
	})
	return
ERR:
	c.JSON(http.StatusOK, gin.H{
		"code": "-1",
		"msg":  msg,
	})
	return
}

func DownloadHandler(c *gin.Context) {
	var err error
	var msg string
	var data []byte
	var f *os.File

	fname := c.Request.FormValue("filehash")
	fm := meta.GetFileMeta(fname)
	f, err = os.Open(fm.Location)
	defer func() {
		if f != nil {
			f.Close()
		}
	}()

	if err != nil {
		msg = fmt.Sprintln("failed to Open:", err)
		goto ERR
	}

	data, err = ioutil.ReadAll(f)
	if err != nil {
		msg = fmt.Sprintln("failed to ReadAll:", err)
		goto ERR
	}
	c.Header("Content-Type", "application/octect-stream")
	c.Header("content-disposition", "attachment;filename=\""+fm.FileName+"\"")
	c.Data(http.StatusOK, "application/octect-stream", data)
	return
ERR:
	c.JSON(http.StatusOK, gin.H{
		"code": "-1",
		"msg":  msg,
	})
	return
}

// post request
func UpdateFileMetaHandler(c *gin.Context) {
	var err error
	var msg string
	var data []byte
	var fm meta.FileMeta

	opType := c.Request.FormValue("op")
	fhash := c.Request.FormValue("filehash")
	fname := c.Request.FormValue("filename")
	if opType != "0" {
		msg = fmt.Sprintln("op is not 0:")
		goto ERR
	}

	fm = meta.GetFileMeta(fhash)
	fm.FileName = fname
	meta.UpdateFileMetas(fm)
	data, err = json.Marshal(fm)
	if err != nil {
		msg = fmt.Sprintln("failed to Marshal:", err)
		goto ERR
	}

	c.Data(http.StatusOK, "application/json", data)
	return
ERR:
	c.JSON(http.StatusOK, gin.H{
		"code": "-1",
		"msg":  msg,
	})
	return
}

func FileDeleteHandle(c *gin.Context) {
	var err error
	var msg string
	fhash := c.Request.FormValue("filehash")
	fLocation := meta.GetFileMeta(fhash).Location
	meta.RemoveFileMeta(fhash)

	err = os.Remove(fLocation)
	if err != nil {
		msg = fmt.Sprintln("failed to Remove:", err)
		goto ERR
	}
	c.JSON(http.StatusOK, gin.H{msg: "OK"})
	return
ERR:
	c.JSON(http.StatusOK, gin.H{
		"code": "-1",
		"msg":  msg,
	})
	return
}
