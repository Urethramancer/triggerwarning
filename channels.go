package main

// Channels are used all over the program
var Channels struct {
	mainquit    chan bool
	janitorquit chan bool
}

func initChannels() {
	// App management
	Channels.mainquit = make(chan bool, 0)
	Channels.janitorquit = make(chan bool, 0)
}
