# Trigger Warning
Pull Git repositories and notify watchers on trigger.

## Who or what is this for?
This is an internal utility I use to keep repositories up to date on a server whenever I check in elsewhere.

## Status
Currently it runs as a web server and accepts HTTP or HTTPS requests to the /trigger endpoint, which takes a `name`parameter to select a repository to update (either internal pull or external command `git pull`).

## Features
Currently users can be added and trigger pulls.

## Unimplemented features
- Add users as watchers of specific repositories.
- Notify watchers.

## Licence
MIT.
