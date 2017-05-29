package main

import (
	"os"

	flags "github.com/jessevdk/go-flags"
)

var opt struct {
	Version VersionCmd `command:"version" description:"Print version and exit." alias:"ver"`
	Config  string     `short:"C" long:"config" description:"Path to configuration file." value-name:"CONFIG"`

	// Commands.
	Server ServerCmd `command:"server" description:"Run the server." alias:"serve" alias:"run"`
	List   ListCmd   `command:"list" description:"List users, triggers and watchers."`
	Add    AddCmd    `command:"add" description:"Add users, triggers and watchers."`
}

type VersionCmd struct {
}

// Execute the version command.
func (cmd *VersionCmd) Execute(args []string) error {
	pr("%s %s", program, Version)
	os.Exit(0)
	return nil
}

// ServerCmd arguments.
type ServerCmd struct {
	Verbose bool `short:"v" long:"verbose" description:"Show detailed information."`
}

// Execute the server command.
// Just lets main() run its course.
func (cmd *ServerCmd) Execute(args []string) error {
	return nil
}

func parseFlags() {
	_, err := flags.Parse(&opt)
	if err != nil {
		os.Exit(0)
	}
}
