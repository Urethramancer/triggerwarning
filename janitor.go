// The janitor cleans up expired tokens.
package main

import (
	"os"
	"time"

	"path/filepath"

	"github.com/fsnotify/fsnotify"
)

func startJanitor() {
	info("Starting janitor with cleanup interval %s.", cfg.CleanupInterval)
	ticker := time.NewTicker(getTime(cfg.CleanupInterval))

	flags := fsnotify.Create | fsnotify.Write | fsnotify.Remove
	userwatcher, err := fsnotify.NewWatcher()
	if err != nil {
		crit("Error starting user watcher: %s", err.Error())
		os.Exit(2)
	}

	triggerwatcher, err := fsnotify.NewWatcher()
	if err != nil {
		crit("Error starting trigger watcher: %s", err.Error())
		os.Exit(2)
	}

	userwatcher.Add(cfg.Users)
	triggerwatcher.Add(cfg.Triggers)
	go func() {
		for {
			select {
			case <-ticker.C:
				//				clearTokens()

			case event := <-userwatcher.Events:
				if flags&event.Op&fsnotify.Create == fsnotify.Create {
					loadUser(event.Name)
					continue
				}
				if flags&event.Op&fsnotify.Write == fsnotify.Write {
					loadUser(event.Name)
					continue
				}
				if flags&event.Op&fsnotify.Remove == fsnotify.Remove {
					unloadUser(filepath.Base(event.Name))
					continue
				}

			case event := <-triggerwatcher.Events:
				if flags&event.Op&fsnotify.Create == fsnotify.Create {
					loadTrigger(event.Name)
					continue
				}
				if flags&event.Op&fsnotify.Write == fsnotify.Write {
					loadTrigger(event.Name)
					continue
				}
				if flags&event.Op&fsnotify.Remove == fsnotify.Remove {
					unloadTrigger(filepath.Base(event.Name))
					continue
				}

			case <-Channels.janitorquit:
				info("Shutting down janitor.")
				ticker.Stop()
				userwatcher.Close()
				return
			}
		}
	}()
}
