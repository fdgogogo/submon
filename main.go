package main

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"github.com/op/go-logging"
	"gopkg.in/alecthomas/kingpin.v2"
	"os"
	"os/user"
	"time"
)

var (
	app = kingpin.New("submon", "").Version("0.1")

	debug = app.Flag("debug", "Enable debug mode.").Bool()
	mono = app.Flag("mono", "Monochrome logging output").Bool()
	verbose = app.Flag("verbose", "Enable verbose mode.").Bool()
	quiet = app.Flag("quiet", "Enable quiet mode.").Bool()

	lang = app.Flag("lang", "language, choice: chn, eng").String()

	download = app.Command("download", "Download subtitle for specific file")
	downloadFile = download.Arg("file", "target file path").Required().String()

	watch = app.Command("watch", "Watch direcotry (and children) for change, download subtitle automatically when new file added.")
	targetDir = watch.Flag("dir", "target dir").Short('d').String()
	//fullScan = watch.Flag("full-scan", "should perform full scan at target dir").Short("f").Default(true).Bool()
	maxRetry = watch.Flag("max-retry", "Max retry before give up").Default("3").Int()
	configFile = watch.Flag("config-file", "config file path").Default("~/.config/submon/config.yaml").String()
	exampleConfig = app.Command("example_config", "show example configuration")
)

var (
	logger = logging.MustGetLogger("submon")
)

const (
	coloredFormat = `%{color}%{time:15:04:05.000} ▶ %{level:.4s} %{color:reset} %{message}`
	monoFormat = `%{time:15:04:05.000} ▶ %{level:.4s} %{message}`
)

var DB *gorm.DB
var AppConfig Config
var command string

func init() {
	// Logging
	var backend logging.Backend
	backend = logging.NewLogBackend(os.Stderr, "", 0)

	if *mono {
		backend = logging.NewBackendFormatter(backend, logging.MustStringFormatter(monoFormat))
	} else {
		backend = logging.NewBackendFormatter(backend, logging.MustStringFormatter(coloredFormat))
	}

	leveledBackend := logging.AddModuleLevel(backend)

	if *verbose {
		leveledBackend.SetLevel(logging.DEBUG, "submon")
	} else if *quiet {
		leveledBackend.SetLevel(logging.WARNING, "submon")
	} else {
		leveledBackend.SetLevel(logging.INFO, "submon")
	}

	logging.SetBackend(leveledBackend)

	// 爬取参数
	command = kingpin.MustParse(app.Parse(os.Args[1:]))

	// 创建项目文件路径
	var err error
	usr, err := user.Current()
	if err != nil {
		logger.Fatal(err)
	}

	home := usr.HomeDir + "/.submon/"
	// TODO: Windows

	err = os.MkdirAll(home, 0755)

	if err != nil {
		logger.Fatal(err)
	}

	// 数据库初始化
	dbPath := home + "db.sqlite"
	DB, err = gorm.Open("sqlite3", dbPath)

	if err != nil {
		panic("failed to connect database")
	}
	if *debug {
		DB.LogMode(true)
	}
	DB.AutoMigrate(&ScannedFile{})
}

func main() {
	app.HelpFlag.Short('h')
	kingpin.Version("0.1")

	AppConfig = ReadConfigFile(*configFile)
	if *lang != "" {
		AppConfig.Lang = *lang
	}

	switch command {

	case download.FullCommand():
	//RequestSubtitle(nil, *downloadFile)

	case watch.FullCommand():
		Watch()

	case exampleConfig.FullCommand():
		PrintDefaultConfig()

	default:
		kingpin.Usage()

	}
}

func Watch() {
	logger.Notice("Submon watcher started")
	PrintDBStat()
	//(scanned int, modified int, notSeen int)
	start := time.Now()
	total, video, modified, new := WalkDir(*targetDir)

	logger.Infof("Found video files:       %d (in %d files)", video, total)
	logger.Infof("New files:               %d", new)
	logger.Infof("Modified files:          %d", modified)
	logger.Infof("Scan completed, took %.2f seconds", time.Since(start).Seconds())
}

func PrintDBStat() {
	var count int
	DB.Model(&ScannedFile{}).Count(&count)
	logger.Infof("Scanned files: %d", count)
}
