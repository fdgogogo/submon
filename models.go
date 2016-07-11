package main

import (
	"github.com/jinzhu/gorm"
	"time"
)

type ScannedFile struct {
	gorm.Model
	ID             string `gorm:"primary_key"`
	FirstSeenAt    time.Time
	LastSeenAt     time.Time
	FileModifiedAt time.Time
	FoundAt        time.Time
	Size           int64
	FoundSubtitle  bool
	FailedTimes    int
}
