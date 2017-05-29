package main

import (
	"errors"
	"net/http"
	"strings"
)

// apiTrigger pulls the latest update for a repository and notifies any watchers.
// TODO: Actually implement watchers.
func apiTrigger(w http.ResponseWriter, r *http.Request) error {
	q := r.URL.Query()
	username := q.Get("user")
	token := q.Get("token")
	username = strings.Trim(username, "\n ")
	user := users.Get(username)
	if user == nil {
		return errors.New("unknown user " + username)
	}

	if user.Token != token {
		return errors.New("token mismatch for user " + username)
	}

	name := q.Get("name")
	name = strings.Trim(name, "\n ")
	trigger := triggermap[name]
	if trigger == nil {
		return errors.New("unknown repository " + name)
	}

	info("Updating %s (%s) from %s", trigger.Name, trigger.Path, trigger.URL)
	err := gitPull(trigger.Path)
	return err
}
