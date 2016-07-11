package main

import (
	"crypto/md5"
	"fmt"
	"github.com/deckarep/golang-set"
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

func WalkDir(dir string) {
	err := filepath.Walk(dir, func(path string, f os.FileInfo, err error) error {
		ext := filepath.Ext(path)
		exists := videoFormatsSet.Contains(ext)
		if exists {
			record := CreateOrUpdateRecord(path, f)
			if !record.FoundSubtitle || record.FileModifiedAt != f.ModTime() {
				record.FileModifiedAt = f.ModTime()
				now := time.Now()
				//if record.LastSeenAt - time.Now()
				record.LastSeenAt = now
				found := RequestSubtitle(path)

				if found {
					record.FoundAt = now
					record.FoundSubtitle = true
				} else {
					record.FailedTimes++
				}
			}
			DB.Save(&record)
		} else {

		}
		return nil
	})
	if err != nil {
		fmt.Println(err)
	}
}

func CreateOrUpdateRecord(path string, f os.FileInfo) (fileRecord ScannedFile) {
	hash := fmt.Sprintf("%x", md5.Sum([]byte(path)))
	DB.First(&fileRecord, path)

	if fileRecord == (ScannedFile{}) {
		fileRecord = ScannedFile{
			ID:             hash,
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