package main

import (
	"encoding/json"
	"fmt"
	"go_code/project14/file-store-server/config"
	"go_code/project14/file-store-server/db"
	"go_code/project14/file-store-server/mq"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
)

func processCallback(msg []byte) bool {
	fmt.Println("processCallback call ...")
	// 读取文件信息
	transData := mq.TransferData{}
	err := json.Unmarshal(msg, &transData)
	if err != nil {
		fmt.Println(err.Error())
		return false
	}
	// 获取文件句柄
	// 将文件合并位小块然后放入指订位置 更新myql
	fileDir, err := os.ReadDir(transData.CurLocation)
	if err != nil {
		fmt.Println("os.ReadDir", err.Error())
		return false
	}

	fileMap := make(map[int]*os.File)
	for i := 1; i <= len(fileDir); i++ {

		filePath := filepath.Join(transData.CurLocation, strconv.Itoa(i))
		
		filed, err := os.Open(filePath)
		if err != nil {
			fmt.Println(err.Error())
			return false
		}
		fileMap[i] = filed
	}
	// 合并分块
	err = os.MkdirAll(transData.DestLocation, 0644)
	if err != nil {
		fmt.Println(err.Error())
		return false
	}
	f, err := os.Create(filepath.Join(transData.DestLocation, transData.FileHash))
	if err != nil {
		fmt.Println(err.Error())
		return false
	}
	defer f.Close()

	for i := 1; i < len(fileMap); i++ {
		partFile := fileMap[i]
		buf, err := ioutil.ReadAll(partFile)
		if err != nil {
			fmt.Println(err.Error())
			return false
		}
		_, err = f.Write(buf)
		if err != nil {
			fmt.Println(err.Error())
			return false
		}

	}
	// 返回结果
	// 更新 MySQL
	err = db.UpdateFileLocation(transData.FileHash, filepath.Join(transData.DestLocation, transData.FileHash))
	if err != nil {
		return false
	}
	return true
}

func main() {
	fmt.Println("queue start...")

	mq.StartConsume(config.TransOSSQueueName,
		"transfer_oss",
		processCallback)
	fmt.Println("queue end...")
}
