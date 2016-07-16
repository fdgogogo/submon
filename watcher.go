package main

import (
	"github.com/howeyc/fsnotify"
	"github.com/mitchellh/go-homedir"
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
		logger.Info("Start watching dir " + path)
		err = watcher.Watch(path)
		if err != nil {
			logger.Fatal(err)
		}
	}

	// Hang so program doesn't exit
	<-done

	/* ... do stuff ... */
	watcher.Close()
}
