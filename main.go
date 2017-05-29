package main

import (
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	parseFlags()
	loadConfig()
	initChannels()
	openLogs()
	defer closeLogs()
	loadUsers(cfg.Users)
	loadTriggers(cfg.Triggers)
	startJanitor()

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		Channels.janitorquit <- true
		Channels.mainquit <- true
	}()

	initWeb()
	<-Channels.mainquit
	info("Quit signal received. Shutting down.")
	go stopServers()
	time.Sleep(time.Millisecond * 500)
	return
}
