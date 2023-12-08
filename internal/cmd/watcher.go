package cmd

import (
	"log"

	"github.com/fsnotify/fsnotify"
)

func watcher(dir string) error {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return err
	}
	defer watcher.Close()

	go func() {
		for {
			select {
			case event, ok := <-watcher.Events:
				if !ok {
					return
				}
				if event.Has(fsnotify.Write) {
					log.Println("modified file:", event.Name)
				}
			case err, ok := <-watcher.Errors:
				if !ok {
					return
				}
				log.Println("error:", err)
			}
		}
	}()

	// Add a path.
	err = watcher.Add("/tmp")
	if err != nil {
		return err
	}

	// Block main goroutine forever.
	<-make(chan struct{})
	return nil
}
