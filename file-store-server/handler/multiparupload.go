package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"go_code/project14/file-store-server/cache/redis"
	"go_code/project14/file-store-server/config"
	"go_code/project14/file-store-server/db"
	"go_code/project14/file-store-server/mq"
	"math"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

type MutilPartUploadInfo struct {
	FileHash   string
	ChunkSize  int
	ChunkCount int
	UploadId   string
}

func MutilPartUploaInitdHandler(c *gin.Context) {
	ctx := context.Background()

	//1 parse paramer
	var err error
	var msg string
	var info MutilPartUploadInfo
	var jsonInfo []byte
	username := c.Request.FormValue("username")
	fileHash := c.Request.FormValue("filehash")
	fileName := c.Request.FormValue("filename")
	fileSize, err := strconv.Atoi(c.Request.FormValue("filesize"))
	if err != nil {
		msg = fmt.Sprintln("invalid param:", err)
		goto ERR

	}

	// 实现秒传
	if fileInfo, err := db.GetFileMeta(fileHash); err == nil {
		if fileInfo.FileName.String != "" {
			if fileInfo.FileSize.Int64 == int64(fileSize) {
				// 触发 秒传功能
				db.OnUserFileUpdateFinishied(username, fileHash, fileName, int64(fileSize))
				c.JSON(http.StatusOK, gin.H{
					"msg":  "OK",
					"code": "0",
				})
				return
			}
		}
	}

	// 3. init multi part upload info
	info = MutilPartUploadInfo{
		FileHash:   fileHash,
		ChunkSize:  5 * 1024 * 1024, //5m
		ChunkCount: int(math.Ceil(float64(fileSize) / (5 * 1024 * 1024))),
		UploadId:   username + fmt.Sprintf("%x", time.Now().Unix()),
	}
	jsonInfo, _ = json.Marshal(info)
	// 2.writer cache
	redis.RedisPool().HSet(ctx, "MP_"+info.UploadId, "filehash", info.FileHash)
	redis.RedisPool().HSet(ctx, "MP_"+info.UploadId, "chunksize", info.ChunkSize)
	redis.RedisPool().HSet(ctx, "MP_"+info.UploadId, "chunkcount", info.ChunkCount)

	c.JSON(http.StatusOK, gin.H{
		"code": "0",
		"msg":  jsonInfo,
	})
	return

ERR:
	c.JSON(http.StatusOK, gin.H{
		"code": "-1",
		"msg":  msg,
	})
	return
}

func UploadPartHandler(c *gin.Context) {

	var err error
	var msg string
	var file *os.File
	var buf = make([]byte, 1024*1024)
	// 1 .parse params

	ctx := context.Background()
	uploadId := c.Request.FormValue("uploadid")
	chunkIndex := c.Request.FormValue("index")

	// 3 . get file to store part content
	fPath := "/data/" + uploadId + "/" + chunkIndex
	os.MkdirAll(path.Dir(fPath), 0744)
	file, err = os.Create(fPath)

	if err != nil {
		msg = fmt.Sprintln("Upload part files failed", err)
		goto ERR
	}

	defer func() {
		if file != nil {
			file.Close()
		}
	}()

	for {
		n, err := c.Request.Body.Read(buf)
		file.Write(buf[:n])
		if err != nil {
			break
		}
	}

	// 4 . update redis cache
	redis.RedisPool().HIncrBy(ctx, "MP_"+uploadId, "chkidx_"+uploadId, 1)
	// 5 . return result to client
	c.JSON(http.StatusOK, gin.H{
		"code": "0",
		"msg":  "OK",
	})
	return

ERR:
	c.JSON(http.StatusOK, gin.H{
		"code": "-1",
		"msg":  msg,
	})
	return
}

func CompleteUplaodHandler(c *gin.Context) {
	// parse params
	var err error
	var msg string
	var data []byte
	var chunkCount int
	var transData mq.TransferData
	var suc bool

	ctx := context.Background()
	username := c.Request.FormValue("username")
	uploadId := c.Request.FormValue("uploadid")
	fileSize, _ := strconv.Atoi(c.Request.FormValue("filesize"))
	fileName := c.Request.FormValue("filename")
	fileHash := c.Request.FormValue("filehash")
	// use uploadId to quert redis data and judge

	oldCount, err := strconv.Atoi(redis.RedisPool().HGet(ctx, "MP_"+uploadId, "chunkcount").Val())

	if err != nil {
		msg = fmt.Sprintln("CompleteUplaodHandler strconv.Atoi", err)
		goto ERR

	}

	chunkCount, err = strconv.Atoi(redis.RedisPool().HGet(ctx, "MP_"+uploadId, "chkidx_"+uploadId).Val())
	if err != nil {
		msg = fmt.Sprintln("CompleteUplaodHandler strconv.Atoi", err)
		goto ERR
	}
	if chunkCount != oldCount {
		msg = fmt.Sprintln("CompleteUplaodHandler chunkCount != old Count")
		goto ERR
	}
	// 合并分块
	// 消息队列中写入信息
	transData = mq.TransferData{
		FileHash:     fileHash,
		CurLocation:  "/data/" + uploadId,
		DestLocation: "/data/udata",
	}
	data, err = json.Marshal(transData)
	if err != nil {
		msg = fmt.Sprintln("CompleteUplaodHandler json.Marshal", err)
		goto ERR
	}
	suc = mq.Publish(config.TransExchangeName, config.TransOSSRountingKey, data)

	if !suc {
		msg = fmt.Sprintln("CompleteUplaodHandler mq.Publish fail")
		goto ERR
	}
	// 更新唯一的文件表和用户表
	db.OnUploadFileFinishied(fileHash, fileName, "", int64(fileSize))
	db.OnUserFileUpdateFinishied(username, fileHash, fileName, int64(fileSize))

	// cache
	redis.RedisPool().Expire(ctx, "MP_"+uploadId, 5)
	// 响应处理结果

	c.JSON(http.StatusOK, gin.H{
		"code": "0",
		"msg":  "OK",
	})
	return
ERR:
	c.JSON(http.StatusOK, gin.H{
		"code": "-1",
		"msg":  msg,
	})
	return
}

// 分块取消
func CancelUploadPartHandler(c *gin.Context) {

	ctx := context.TODO()
	// 通知 上传函数停止上传

	username := c.Request.FormValue("username")
	uploadId := c.Request.FormValue("uploadid")
	fileHash := c.Request.FormValue("filehash")
	fileName := c.Request.FormValue("filename")
	// 通知各个线程 结束上传
	// 删除已经存在的分块信息
	fPath := "/data/" + uploadId
	os.RemoveAll(fPath)

	// 删除redis信息
	redis.RedisPool().Expire(ctx, "MP_"+uploadId, 5)
	// 更新MySQL记录
	db.OnUserFileCancel(username, fileHash, fileName)
	c.JSON(http.StatusOK, gin.H{
		"code": "0",
		"msg":  "OK",
	})
	return
}

func MutilPartUploadStatusHandler(c *gin.Context) {
	// 检查分块上传状态信息是否有效
	ctx := context.Background()
	uploadId := c.Request.FormValue("uploadid")

	// 判断是否存在
	fPath := filepath.Join("/data", uploadId)
	if _, err := os.Stat(fPath); os.IsNotExist(err) {
		c.JSON(http.StatusOK, gin.H{
			"code": "-1",
			"msg":  "upload file not exist",
		})
		return
	}
	// 获取分块上传信息
	total, _ := strconv.Atoi(redis.RedisPool().HGet(ctx, "MP_"+uploadId, "chunkCount").Val())

	// 获取已经上传的分块信息
	curTotal, _ := strconv.Atoi(redis.RedisPool().HGet(ctx, "MP_"+uploadId, "chkidx_"+uploadId).Val())
	// 完成率
	res := float64(curTotal) / float64(total)
	c.JSON(http.StatusOK, gin.H{
		"code": "0",
		"msg":  res,
	})
	return
}

// // 断点续传

// func BreakpointContinuation(w http.ResponseWriter, r http.Request) {

// }
