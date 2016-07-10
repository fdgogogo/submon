package main

import (
	"os"
	"math"
	"fmt"
	"crypto/md5"
)

func ComputeFileHash(filePath string) (hash string) {
	fp, err := os.Open(filePath)
	if err != nil {

	}
	stat, err := fp.Stat()
	if err != nil {

	}
	size := float64(stat.Size())
	sample_positions := [4]int64{
		4 * 1024,
		int64(math.Floor(size / 3 * 2)),
		int64(math.Floor(size / 3)),
		int64(size - 8 * 1024)}
	var samples [4][]byte
	for i, position := range sample_positions {
		samples[i] = make([]byte, 4 * 1024)
		fp.ReadAt(samples[i], position)
	}
	for _, sample := range samples {
		if len(hash) > 0 {
			hash += ";"
		}
		hash += fmt.Sprintf("%x", md5.Sum(sample))
	}

	return
}
