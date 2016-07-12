package main

import (
	"bytes"
	"encoding/json"
	"io"
	"io/ioutil"
	"mime"
	"net/http"
	"net/url"
	"os"
	"path"
	"time"
)

type ScannedFile struct {
	// Model
	ID             string `gorm:"type:char(32);primary_key"`
	Path           string
	FirstSeenAt    time.Time
	LastSeenAt     time.Time
	FileModifiedAt time.Time
	FoundAt        time.Time
	Size           int64
	FoundSubtitle  bool
	FailedTimes    int
}

type FileInfo struct {
	Ext  string // 文件扩展名
	Link string // 文件下载链接
}

type SubInfo struct {
	Desc  string     // 备注信息
	Delay int32      // 字幕相对于视频的延迟时间，单位是毫秒
	Files []FileInfo // 包含文件信息的Array。 注：一个字幕可能会包含多个字幕文件，例如：idx+sub格式
}

func (self *ScannedFile) Save() {
	DB.Save(&self)
}

func (self *ScannedFile) RequestSubtitle() {
	var found bool

	logger.Notice("Start searching subtitles for " + path.Base(self.Path))
	Url, err := url.Parse("https://www.shooter.cn/api/subapi.php")

	if err != nil {
		logger.Error(err)
	}

	hash := ComputeFileHash(self.Path)

	parameters := url.Values{}
	parameters.Add("filehash", hash)
	parameters.Add("pathinfo", self.Path)
	parameters.Add("format", "json")
	parameters.Add("lang", AppConfig.Lang)
	Url.RawQuery = parameters.Encode()

	req, err := http.NewRequest("POST", Url.String(), bytes.NewBuffer(make([]byte, 0)))
	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	logger.Debug("Response status:", resp.Status)
	body, _ := ioutil.ReadAll(resp.Body)

	var data []SubInfo
	json.Unmarshal(body, &data)
	for _, subinfo := range data {
		for _, fileinfo := range subinfo.Files {
			DownloadSubtitle(self.Path, fileinfo.Link)
		}
		found = true
	}

	if found {
		self.FoundAt = time.Now()
		self.FoundSubtitle = true
	} else {
		self.FailedTimes++
	}
	self.Save()
}

func DownloadSubtitle(filePath string, link string) {
	response, req_err := http.Get(link)
	if req_err != nil {
		logger.Error(req_err)
	}
	defer response.Body.Close()

	_, params, _ := mime.ParseMediaType(response.Header["Content-Disposition"][0])
	filename := params["filename"] // set to "foo.png"

	target := path.Dir(filePath) + "/" + filename
	out, os_err := os.Create(target)
	if os_err != nil {
		logger.Error(os_err)
	} else {
		_, dl_err := io.Copy(out, response.Body)
		if dl_err != nil {
			logger.Error(dl_err)
		}
	}
	defer out.Close()
	logger.Info("Downloaded " + target)

	return
}
