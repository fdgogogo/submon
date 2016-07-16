package main

import (
	"github.com/howeyc/fsnotify"
	"github.com/mitchellh/go-homedir"
	"path/filepath"
	"os"
)

func Watch() {
	watcher, err := fsnotify.NewWatcher()

	if err != nil {
		logger.Fatal(err)
	}

	done := make(chan bool)

	// Process events
	go func() {
		for {
			select {
			case ev := <-watcher.Event:
				if ev.IsCreate() {
					stat, err := os.Stat(ev.Name)
					if err != nil {
						logger.Error(err)
					}
					if stat.IsDir() {
						logger.Debug("Start watching dir " + ev.Name)
						watcher.Watch(ev.Name)
					}
				}

				logger.Debug("event:", ev)
			case err := <-watcher.Error:
				logger.Debug("error:", err)
			}
		}
	}()

	for _, wc := range AppConfig.Watch {
		path, err := homedir.Expand(wc.Path)
		if err != nil {
			logger.Fatal(err)
		}
		if wc.Recursive {
			logger.Debug("Start watching dir " + path + "(Recursive)")
			paths := ListSubDir(path)
			for _, p := range paths {
				err = watcher.Watch(p)
			}
		} else {
			logger.Debug("Start watching dir " + path)
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
