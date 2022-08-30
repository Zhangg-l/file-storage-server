package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"strconv"

	jsonit "github.com/json-iterator/go"
)

func main() {
	username := "admin"
	token := "dc9b7cda5edf0df6491f6a7e3c33f33116617831"
	filesha1 := "809017202e3968daf4bd68d37022dbc7f4baae8d"
	filename := "/home/vscode-server-linux-x64.tar.gz"
	// fileSize := 0
	// 默契post请求
	resp, err := http.PostForm(
		"http://localhost:8888/file/mpupload/init",
		url.Values{
			"username": {username},
			"token":    {token},
			"filehash": {filesha1},
			"filesize": {"55068508"},
		})
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(-1)
	}
	defer resp.Body.Close()
	bodyDate, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(-1)
	}
	uploadId := jsonit.Get(bodyDate, "data").Get("UploadId").ToString()
	chunkSize := jsonit.Get(bodyDate, "data").Get("ChunkSize").ToInt()

	fmt.Printf("uploadid: %s  chunksize: %d\n", uploadId, chunkSize)

	//
	tURL := "http://localhost:8888/file/mpupload/uppart?" +
		"username=admin&token=" + token + "&uploadid=" + uploadId

	multipartUpload(filename, tURL, chunkSize)
	// 请求完成

	resp, err = http.PostForm("http://localhost:8888/file/mpupload/complete",
		url.Values{
			"username": {username},
			"token":    {token},
			"filehash": {filesha1},
			"filename": {filename},
			"filesize": {"23233372"},
			"uploadid": {uploadId}})
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(-1)
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(-1)
	}
	fmt.Printf("complete result: %s\n", string(body))
}

func multipartUpload(filename string, tURL string, chunkSize int) error {
	f, err := os.Open(filename)
	if err != nil {
		fmt.Println(err)
		return err
	}

	defer f.Close()

	bfRD := bufio.NewReader(f)

	index := 0

	ch := make(chan int)
	buf := make([]byte, chunkSize)

	for {
		n, err := bfRD.Read(buf)
		if n <= 0 {
			break
		}
		index++
		bufCopied := make([]byte, chunkSize)
		copy(bufCopied, buf)
		go func(b []byte, curIdx int) {
			fmt.Printf("upload_size: %d\n", len(b))
			resp, err := http.Post(tURL+"&index="+strconv.Itoa(curIdx),
				"multipart.form-data",
				bytes.NewBuffer(b))

			if err != nil {
				fmt.Println(err)
			}
			body, er := ioutil.ReadAll(resp.Body)

			fmt.Printf("%+v %+v\n", string(body), er)
			resp.Body.Close()

			ch <- curIdx
		}(bufCopied[:n], index)

		if err != nil {
			if err == io.EOF {
				break
			} else {
				fmt.Println(err.Error())
			}
		}

	}

	for idx := 0; idx < index; idx++ {
		select {
		case res := <-ch:
			fmt.Println(res)
		}
	}
	return nil
}
