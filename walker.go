package main

import (
	"crypto/md5"
	"fmt"
	"github.com/deckarep/golang-set"
	"github.com/mitchellh/go-homedir"
	"os"
	"path/filepath"
	"time"
)

var videoFormats = []interface{}{
	".webm",
	".mkv",
	".flv",
	".flv",
	".vob",
	".ogv",
	".ogg",
	".drc",
	".gifv",
	".mng",
	".avi",
	".mov",
	".qt",
	".wmv",
	".yuv",
	".rm",
	".rmvb",
	".asf",
	".amv",
	".mp4",
	".m4p",
	".m4v",
	".mpg",
	".mp2",
	".mpeg",
	".mpe",
	".mpv",
	".mpg",
	".mpeg",
	".m2v",
	".m4v",
	".svi",
	".3gp",
	".3g2",
	".mxf",
	".roq",
	".nsv",
	".flv",
	".f4v",
	".f4p",
	".f4a",
	".f4b",
}

var videoFormatsSet = mapset.NewSetFromSlice(videoFormats)

func IsVideoFile(p string) (b bool) {
	ext := filepath.Ext(p)
	return videoFormatsSet.Contains(ext)
}

func WalkDir(dir string) (total int, video int, modified int, new int) {
	taskQueue := NewTaskQueue().Run()
	logger.Debugf("%s", taskQueue)
	absDir, err := homedir.Expand(dir)
	if err != nil {
		logger.Fatal(err)
	}

	err = filepath.Walk(absDir, func(path string, f os.FileInfo, err error) error {
		total++
		isVideoFile := IsVideoFile(path)

		if !isVideoFile {
			return nil
		}
		video++

		record, shouldRequest := CreateOrUpdateRecord(path, f)
		if shouldRequest {
			if DB.NewRecord(record) {
				new++
			} else {
				modified++
			}
			taskQueue <- &record
		}
		return nil

	})

	if err != nil {
		logger.Error(err)
	}
	close(taskQueue)

	return
}

func CreateOrUpdateRecord(path string, f os.FileInfo) (fileRecord VideoFile, shouldRequest bool) {
	hash := fmt.Sprintf("%x", md5.Sum([]byte(path)))
	now := time.Now()
	shouldRequest = true

	DB.Where("id = ?", hash).First(&fileRecord)

	if fileRecord == (VideoFile{}) {
		fileRecord = VideoFile{
			ID:             hash,
			Path:           path,
			FirstSeenAt:    time.Now(),
			FileModifiedAt: f.ModTime(),
			Size:           f.Size(),
			FoundSubtitle:  false,
			FailedTimes:    0,
		}
		DB.Create(&fileRecord)
	}

	fileRecord.LastSeenAt = now
	fileRecord.FileModifiedAt = f.ModTime()

	if fileRecord.FoundSubtitle && fileRecord.FileModifiedAt == f.ModTime() {
		// 忽略已经找到字幕的且未修改的文件
		logger.Debugf("%s already has subtitle, skipping", f.Name())
		shouldRequest = false
	}

	if fileRecord.FailedTimes >= AppConfig.Watch.MaxRetry {
		// 忽略重试超过一定次数的文件
		logger.Debugf("%s failed too many times, skipping", f.Name())
		shouldRequest = false
	}

	fileRecord.Save()

	return
}
