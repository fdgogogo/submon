package main

import ()
import (
	"gopkg.in/alecthomas/kingpin.v2"
	"os"
)

var (
	app = kingpin.New("shooter-subtitle-worker", "").Version("0.1")

	download = app.Command("download", "Download subtitle for specific file")
	download_file = download.Arg("file", "target file path").Required().String()
	download_lang = download.Flag("lang", "language, choice: chn, eng").Default("chn").String()

	watch = app.Command("watch", "Watch direcotry (and children) for change, download subtitle automatically when new file added.")
	watch_dir = watch.Flag("dir", "target dir").Short('d').String()
	config_file = watch.Flag("config-file", "config file path").String()
	example_config = app.Command("example_config", "show example configuration")
)

func main() {
	app.HelpFlag.Short('h')
	kingpin.Version("0.1")

	parsed := kingpin.MustParse(app.Parse(os.Args[1:]))

	switch parsed {

	case download.FullCommand():
		RequestSubtitle(*download_file, *download_lang)

	case watch.FullCommand():
		ReadConfig("./config.yaml")
		WalkDir(*watch_dir)

	case example_config.FullCommand():
		PrintDefaultConfig()

	default:
		kingpin.Usage()

	}
}
