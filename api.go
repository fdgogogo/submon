package main

import (
	"bytes"
	"encoding/json"
	"io"
	"io/ioutil"
	"log"
	"mime"
	"net/http"
	"net/url"
	"os"
	"path"
)

type Fileinfo struct {
	Ext  string // 文件扩展名
	Link string // 文件下载链接
}

type Subinfo struct {
	Desc  string     // 备注信息
	Delay int32      // 字幕相对于视频的延迟时间，单位是毫秒
	Files []Fileinfo // 包含文件信息的Array。 注：一个字幕可能会包含多个字幕文件，例如：idx+sub格式
}

var logger = log.New(os.Stdout, "", 0)
var err_logger = log.New(os.Stderr, "", 0)

func RequestSubtitle(filePath string, lang string) {
	logger.Println("Start searching subtitles for " + path.Base(filePath))
	Url, err := url.Parse("https://www.shooter.cn/api/subapi.php")

	if err != nil {
		err_logger.Println(err)
	}

	hash := ComputeFileHash(filePath)

	parameters := url.Values{}
	parameters.Add("filehash", hash)
	parameters.Add("pathinfo", filePath)
	parameters.Add("format", "json")
	parameters.Add("lang", lang)
	Url.RawQuery = parameters.Encode()

	req, err := http.NewRequest("POST", Url.String(), bytes.NewBuffer(make([]byte, 0)))
	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	logger.Println("Response status:", resp.Status)
	body, _ := ioutil.ReadAll(resp.Body)

	var data []Subinfo
	json.Unmarshal(body, &data)
	for _, subinfo := range data {
		for _, fileinfo := range subinfo.Files {
			DownloadSubtitle(filePath, fileinfo.Link)
		}
	}
}

func DownloadSubtitle(filePath string, link string) {
	response, req_err := http.Get(link)
	if req_err != nil {
		err_logger.Println(req_err)
	}
	defer response.Body.Close()

	_, params, _ := mime.ParseMediaType(response.Header["Content-Disposition"][0])
	filename := params["filename"] // set to "foo.png"

	target := path.Dir(filePath) + "/" + filename
	out, os_err := os.Create(target)
	if os_err != nil {
		err_logger.Println(os_err)
	} else {
		_, dl_err := io.Copy(out, response.Body)
		if dl_err != nil {
			err_logger.Println(dl_err)
		}
	}
	defer out.Close()
	logger.Println("Downloaded " + target)

	return
}
