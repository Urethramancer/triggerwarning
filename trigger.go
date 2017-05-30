package main

import (
	"io/ioutil"
	"os"

	"path/filepath"
)

var triggermap = make(map[string]*Trigger, 0)

// Trigger holds details for a Git repository to update when triggered.
type Trigger struct {
	// Name is used to set notification subscriptions for watchers.
	Name string `required:"true" positional-arg-name:"name" description:"One-word name of the trigger." json:"name"`
	// Repository is the URL of the repository to clone and update. This is only used for the initial clone operation.
	URL string `required:"true" positional-arg-name:"url" description:"Repository URL to handle." json:"url"`
	// Path is where the cloned repository will be stored and updated.
	Path string `required:"true" positional-arg-name:"path" description:"Directory to clone into." json:"path"`
	// CLI makes use of the system's command line git command instead of go-git.
	CLI bool `json:"cli,omitempty"`
}

func createTrigger(name, url, path string, cli bool) error {
	// if cross.Exists(path) {
	// 	return errors.New("repository path already exists")
	// }

	var err error
	if cli {
		err = gitExternalClone(url, path)
	} else {
		err = gitClone(url, path)
	}

	if err != nil {
		return err
	}

	t := Trigger{
		Name: name,
		Path: path,
		URL:  url,
		CLI:  cli,
	}

	name = filepath.Join(cfg.Triggers, name)
	SaveJSON(name, &t)
	return nil
}

func loadTriggers(path string) {
	triggers, err := ioutil.ReadDir(path)
	info("Loading triggers from from %s", path)
	if err != nil {
		crit("Error loading triggers: %s", err.Error())
		os.Exit(2)
	}

	for _, t := range triggers {
		if !t.IsDir() {
			fn := filepath.Join(path, t.Name())
			loadTrigger(fn)
		}
	}
	info("Loaded %d trigger(s).", len(triggermap))
}

func loadTrigger(fn string) {
	var trigger Trigger
	LoadJSON(fn, &trigger)
	if trigger.Name != "" {
		triggermap[trigger.Name] = &trigger
		info("Loaded trigger %s", trigger.Name)
	}
}

func unloadTrigger(name string) {
	_, ok := triggermap[name]
	if ok {
		delete(triggermap, name)
		info("Removed trigger %s", name)
	}
}
