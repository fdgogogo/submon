package main

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"github.com/mitchellh/go-homedir"
	"github.com/op/go-logging"
	"gopkg.in/alecthomas/kingpin.v2"
	"os"
	"time"
)

var (
	app = kingpin.New("submon", "").Version("0.1")

	debug   = app.Flag("debug", "Enable debug mode.").Bool()
	mono    = app.Flag("mono", "Monochrome logging output").Bool()
	verbose = app.Flag("verbose", "Enable verbose mode.").Bool()
	quiet   = app.Flag("quiet", "Enable quiet mode.").Bool()

	lang = app.Flag("lang", "language, choice: chn, eng").String()

	download     = app.Command("download", "Download subtitle for specific file")
	downloadFile = download.Arg("file", "target file path").Required().String()

	watch         = app.Command("watch", "Watch direcotry (and children) for change, download subtitle automatically when new file added.")
	targetDir     = watch.Flag("dir", "target dir").Short('d').String()
	noFullScan    = watch.Flag("--no-full-scan", "should perform full scan at target dir").Bool()
	maxRetry      = watch.Flag("max-retry", "Max retry before give up").Default("3").Int()
	configFile    = watch.Flag("config-file", "config file path, default to ~/.submon/config.yaml").Default("~/.submon/config.yaml").String()
	exampleConfig = app.Command("example_config", "show example configuration")
)

var (
	logger = logging.MustGetLogger("submon")
)

const (
	coloredFormat = `%{color}%{time:15:04:05.000} ▶ %{level:.4s} %{color:reset} %{message}`
	monoFormat    = `%{time:15:04:05.000} ▶ %{level:.4s} %{message}`
)

var DB *gorm.DB
var AppConfig Config
var command string

func main() {
	app.HelpFlag.Short('h')
	kingpin.Version("0.1")
	configureApp()

	if *lang != "" {
		AppConfig.Lang = *lang
	}

	switch command {

	case download.FullCommand():
	//RequestSubtitle(nil, *downloadFile)

	case watch.FullCommand():
		WatchCommand()

	case exampleConfig.FullCommand():
		PrintDefaultConfig()

	default:
		kingpin.Usage()

	}
}
func configureApp() {
	var backend logging.Backend

	command = kingpin.MustParse(app.Parse(os.Args[1:]))

	AppConfig = ReadConfigFile(*configFile)
	AppConfig = UpdateAppConfig(AppConfig)

	// Logging
	backend = logging.NewLogBackend(os.Stderr, "", 0)
	switch AppConfig.LogFormat {
	case "mono":
		backend = logging.NewBackendFormatter(backend, logging.MustStringFormatter(monoFormat))
	case "color":
		backend = logging.NewBackendFormatter(backend, logging.MustStringFormatter(coloredFormat))
	}
	leveledBackend := logging.AddModuleLevel(backend)
	leveledBackend.SetLevel(logging.GetLevel(AppConfig.LogLevel), "submon")
	logging.SetBackend(leveledBackend)

	// 创建项目文件路径
	var err error
	usrHome, err := homedir.Dir()
	home := usrHome + "/.submon/"
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
	DB.AutoMigrate(&VideoFile{})
}

// command-line flag override
func UpdateAppConfig(in Config) (out Config) {
	out = in
	if *lang != "" {
		out.Lang = *lang
	}

	if *verbose {
		out.LogLevel = "DBEUG"
	}

	if *quiet {
		out.LogLevel = "NOTICE"
	}

	if *mono {
		out.LogFormat = "mono"
	}

	out.Debug = *debug

	if *targetDir != "" {
		wc := WatchConfig{Path: *targetDir, NoFullScan: *noFullScan, MaxRetry: *maxRetry}
		out.Watch = []WatchConfig{wc}
	}

	if *maxRetry != 0 {
		for _, w := range out.Watch {
			w.MaxRetry = *maxRetry
		}
	}

	if *noFullScan {
		for _, w := range out.Watch {
			w.NoFullScan = true
		}
	}

	return
}

func WatchCommand() {
	logger.Notice("Submon watcher started")
	PrintDBStat()
	//(scanned int, modified int, notSeen int)
	start := time.Now()
	for _, wc := range AppConfig.Watch {

		total, video, modified, new := WalkDir(wc.Path)
		logger.Infof("Found video files:       %d (in %d files)", video, total)
		logger.Infof("New files:               %d", new)
		logger.Infof("Modified files:          %d", modified)
		logger.Infof("Scan completed, took %.2f seconds", time.Since(start).Seconds())
	}
	//Watch()
}

func PrintDBStat() {
	var count int
	DB.Model(&VideoFile{}).Count(&count)
	logger.Infof("Scanned files: %d", count)
}
