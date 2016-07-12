package main

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"gopkg.in/alecthomas/kingpin.v2"
	"log"
	"os"
)

var (
	app  = kingpin.New("submon", "").Version("0.1")
	lang = app.Flag("lang", "language, choice: chn, eng").String()

	download     = app.Command("download", "Download subtitle for specific file")
	downloadFile = download.Arg("file", "target file path").Required().String()

	watch         = app.Command("watch", "Watch direcotry (and children) for change, download subtitle automatically when new file added.")
	watchDir      = watch.Flag("dir", "target dir").Short('d').String()
	configFile    = watch.Flag("config-file", "config file path").Default("~/.config/submon/config.yaml").String()
	exampleConfig = app.Command("example_config", "show example configuration")
)

var (
	logger     = log.New(os.Stdout, "", 0)
	err_logger = log.New(os.Stderr, "", 0)
)

var DB *gorm.DB
var AppConfig Config

func init() {
	var err error
	home := "~/.submon/"
	os.MkdirAll(home, 0666)
	DB, err = gorm.Open("sqlite3", home + "db.sqlite")
	if err != nil {
		panic("failed to connect database")
	}
	DB.AutoMigrate(&ScannedFile{})

}

func main() {
	app.HelpFlag.Short('h')
	kingpin.Version("0.1")

	parsed := kingpin.MustParse(app.Parse(os.Args[1:]))

	AppConfig = ReadConfig(*configFile)
	if *lang != "" {
		AppConfig.Lang = *lang
	}

	switch parsed {

	case download.FullCommand():
		RequestSubtitle(*downloadFile)

	case watch.FullCommand():
		WalkDir(*watchDir)

	case exampleConfig.FullCommand():
		PrintDefaultConfig()

	default:
		kingpin.Usage()

	}
}
