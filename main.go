package main

import ()
import "gopkg.in/alecthomas/kingpin.v2"

var (
	file = kingpin.Flag("file", "filename").String()
)

func main() {
	kingpin.Version("0.1")
	kingpin.Parse()
	RequestSubtitle(*file, "chn")
}
