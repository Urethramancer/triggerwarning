# Trigger Warning [![Build Status](https://travis-ci.org/Urethramancer/triggerwarning.svg?branch=master)](https://travis-ci.org/Urethramancer/triggerwarning)
Pull Git repositories and notify watchers on trigger.

## Who or what is this for?
This is an internal utility I use to keep repositories up to date on a server whenever I check in elsewhere.

## Requirements
- github.com/Urethramancer/cross
- github.com/Urethramancer/slog
- github.com/fsnotify/fsnotify
- github.com/gorilla/mux
- github.com/jessevdk/go-flags
- golang.org/x/crypto/sha3
- golang.org/x/crypto/ssh
- gopkg.in/src-d/go-git.v4
- gopkg.in/src-d/go-git.v4/plumbing/transport/ssh

These should be fetched automatically.

## Status
Currently it runs as a web server and accepts HTTP or HTTPS requests to the /trigger endpoint, which takes a `name`parameter to select a repository to update (either internal pull or external command `git pull`).

## Features
Currently users can be added and trigger pulls.

## Unimplemented features
- Add users as watchers of specific repositories.
- Notify watchers.

## Licence
MIT.
