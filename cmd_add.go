package main

import "os"

// AddCmd arguments.
type AddCmd struct {
	User    AddUser    `command:"user" description:"Add a user and generate a password."`
	Trigger AddTrigger `command:"trigger" description:"Add a trigger for a new repository."`
	Watcher AddWatcher `command:"watcher" description:"Add a user as a watcher for a trigger."`
}

// AddUser arguments.
type AddUser struct {
	Args struct {
		Name  string `required:"true" positional-arg-name:"username" description:"The new user's login name."`
		Email string `required:"true" positional-arg-name:"e-mail" description:"E-mail address for notifications."`
	} `positional-args:"true"`
}

// Execute the add user command.
func (cmd *AddUser) Execute(args []string) error {
	loadConfig()
	t, err := createUser(cmd.Args.Name, cmd.Args.Email)
	if err != nil {
		pr("Error creating user: %s", err.Error())
		os.Exit(2)
	}

	pr("API token: %s", t)
	os.Exit(0)
	return nil
}

// AddTrigger arguments.
type AddTrigger struct {
	CLI     bool    `short:"c" long:"command" description:"Use external git command to clone and update."`
	Trigger Trigger `positional-args:"true"`
}

// Execute the add trigger command.
func (cmd *AddTrigger) Execute(args []string) error {
	loadConfig()

	err := createTrigger(cmd.Trigger.Name, cmd.Trigger.URL, cmd.Trigger.Path, cmd.CLI)

	if err != nil {
		crit("Couldn't create trigger: %s", err.Error())
		os.Exit(1)
	}

	os.Exit(0)
	return nil
}

// AddWatcher arguments.
type AddWatcher struct {
}

// Execute the add trigger command.
func (cmd *AddWatcher) Execute(args []string) error {
	loadConfig()
	os.Exit(0)
	return nil
}
