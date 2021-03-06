package main

import (
	"github.com/howeyc/fsnotify"
	"github.com/mitchellh/go-homedir"
	"os"
	"path/filepath"
)

func Watch() {
	watcher, err := fsnotify.NewWatcher()

	if err != nil {
		logger.Fatal(err)
	}

	done := make(chan bool)
	taskQueue := NewTaskQueue().Run()

	// Process events
	go func() {
		for {
			select {
			case ev := <-watcher.Event:
				EventFunc(ev, watcher, taskQueue)
			//logger.Debug("event:", ev)
			case err := <-watcher.Error:
				logger.Error("error:", err)
			}
		}
	}()

	for _, wc := range AppConfig.Watch.Dirs {
		path, err := homedir.Expand(wc.Path)
		if err != nil {
			logger.Fatal(err)
		}
		if wc.Recursive {
			logger.Notice("Start watching dir " + path + "(Recursive)")
			paths := ListSubDir(path)
			for _, p := range paths {
				err = watcher.Watch(p)
			}
		} else {
			logger.Notice("Start watching dir " + path)
			err = watcher.Watch(path)

		}
		if err != nil {
			logger.Fatal(err)
		}
	}

	// Hang so program doesn't exit
	<-done

	/* ... do stuff ... */
	watcher.Close()
}

func EventFunc(ev *fsnotify.FileEvent, watcher *fsnotify.Watcher, tasksQueue chan *VideoFile) {
	if ev.IsCreate() {
		stat, err := os.Stat(ev.Name)
		if err != nil {
			logger.Error(err)
		}
		if stat.IsDir() {
			logger.Debug("Start watching dir " + ev.Name)
			watcher.Watch(ev.Name)
		} else if IsVideoFile(ev.Name) {
			if ev.IsCreate() {
				logger.Info("Found new video file:", ev.Name)
			}
			record, shouldRequest := CreateOrUpdateRecord(ev.Name, stat)
			if shouldRequest {
				record.RequestSubtitle()
			}
		}
	}
}

func ListSubDir(dir string) (dirs []string) {
	err := filepath.Walk(dir, func(path string, f os.FileInfo, err error) error {
		if f.IsDir() {
			dirs = append(dirs, path)
		}
		return nil
	})

	if err != nil {
		logger.Error(err)
	}
	return
}
