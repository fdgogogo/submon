package main

import (
	"time"
)

type ScannedFile struct {
	ID             string `gorm:"primary_key"`
	FirstSeenAt    time.Time
	LastSeenAt     time.Time
	FileModifiedAt time.Time
	FoundAt        time.Time
	Size           int64
	FoundSubtitle  bool
	FailedTimes    int
}
