package main

import "os"

// ListCmd arguments.
type ListCmd struct {
	Verbose  bool        `short:"v" long:"verbose" description:"Show detailed information."`
	Triggers TriggerList `command:"triggers" description:"List triggers." alias:"t"`
}

type TriggerList struct {
}

// Execute the list command.
// Does its thing and exits before main().
func (cmd *TriggerList) Execute(args []string) error {
	pr("Listing triggers.")
	os.Exit(0)
	return nil
}
