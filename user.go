package main

import (
	"errors"
	"io/ioutil"
	"os"
	"path/filepath"
	"sync"
)

type Users struct {
	sync.RWMutex
	list   map[string]*User
	tokens map[string]*User
}

func (u *Users) Add(user *User) {
	u.Lock()
	defer u.Unlock()
	u.list[user.Username] = user
	u.tokens[user.Token] = user
}

func (u *Users) Get(name string) *User {
	u.RLock()
	defer u.RUnlock()
	user, ok := u.list[name]
	if !ok {
		return nil
	}

	return user
}

func (u *Users) Delete(name string) {
	u.Lock()
	defer u.Unlock()
	user, ok := u.list[name]
	if !ok {
		return
	}

	delete(users.list, name)
	delete(users.tokens, user.Token)
}

var users Users

type User struct {
	// Username is mostly used as a reference by the admin and internal systems.
	Username string `json:"username"`
	// Token is a permanent API token to allow access to triggers and add users to watchers.
	Token string `json:"token"`
	// Email is where notifications for watchers are sent.
	Email string `json:"email"`
}

func init() {
	users.Lock()
	users.list = make(map[string]*User)
	users.tokens = make(map[string]*User)
	users.Unlock()
}

// createUser returns the generated API token.
func createUser(name, email string) (string, error) {
	user := users.Get(name)
	if user != nil {
		return "", errors.New("user already exists")
	}

	hash := hashString(genString(32, true))
	user = &User{
		Username: name,
		Token:    hash,
		Email:    email,
	}

	fn := filepath.Join(cfg.Users, name)
	SaveJSON(fn, user)
	return hash, nil
}

func loadUsers(path string) {
	userdir, err := ioutil.ReadDir(path)
	info("Loading users from %s", path)
	if err != nil {
		crit("Error loading users: %s", err.Error())
		os.Exit(2)
	}

	for _, u := range userdir {
		if !u.IsDir() {
			fn := filepath.Join(path, u.Name())
			loadUser(fn)
		}
	}
	info("Loaded %d user(s).", len(users.list))
}

func loadUser(fn string) {
	var user User
	LoadJSON(fn, &user)
	if user.Username != "" {
		users.Add(&user)
		info("Loaded user %s", user.Username)
	}
}

func unloadUser(name string) {
	users.Delete(name)
	info("Removed user %s", name)
}
