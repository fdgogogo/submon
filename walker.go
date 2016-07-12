package main

import (
	"crypto/md5"
	"fmt"
	"github.com/deckarep/golang-set"
	"os"
	"path/filepath"
	"sync"
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

func WalkDir(dir string) (total int, video int, modified int, new int) {

	const workers = 4
	tasks := make(chan *ScannedFile)
	var wg sync.WaitGroup
	for i := 0; i < workers; i++ {
		wg.Add(1)
		go func() {
			for task := range tasks {
				task.RequestSubtitle()
			}
			wg.Done()
		}()
	}

	err := filepath.Walk(dir, func(path string, f os.FileInfo, err error) error {
		total++
		ext := filepath.Ext(path)
		isVideoFile := videoFormatsSet.Contains(ext)

		if !isVideoFile {
			return nil
		}
		video++

		now := time.Now()
		record := CreateOrUpdateRecord(path, f)
		record.LastSeenAt = now
		record.FileModifiedAt = f.ModTime()

		if !(!record.FoundSubtitle || record.FileModifiedAt != f.ModTime()) {
			// 忽略已经找到字幕的文件且未修改的文件
			record.Save()
			return nil
		}

		if DB.NewRecord(record) {
			new++
		} else {
			modified++
		}
		tasks <- &record
		return nil

	})
	if err != nil {
		logger.Error(err)
	}

	close(tasks)
	wg.Wait()

	return
}

func CreateOrUpdateRecord(path string, f os.FileInfo) (fileRecord ScannedFile) {
	hash := fmt.Sprintf("%x", md5.Sum([]byte(path)))
	DB.Where("id = ?", hash).First(&fileRecord)

	if fileRecord == (ScannedFile{}) {
		fileRecord = ScannedFile{
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
	return
}
