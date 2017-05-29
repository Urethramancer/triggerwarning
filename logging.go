package main

import (
	"os"

	"github.com/Urethramancer/slog"
)

var msgFile, errorFile *os.File

func openLogFile(fn string) *os.File {
	f, err := os.OpenFile(fn, os.O_APPEND|os.O_CREATE|os.O_RDWR, 0600)
	if err != nil {
		crit("Error opening file: %s", err.Error())
		os.Exit(2)
	}
	return f
}

func openLogs() {
	if cfg.Logs.Info != "" {
		msgFile = openLogFile(cfg.Logs.Info)
		slog.SetMsgFile(msgFile)
		slog.SetWarnFile(msgFile)
	}

	if cfg.Logs.Error != "" {
		errorFile = openLogFile(cfg.Logs.Error)
		slog.SetErrorFile(errorFile)
	}
}

func closeLogs() {
	if msgFile != nil {
		slog.SetMsgFile(os.Stdout)
		slog.SetWarnFile(os.Stdout)
		err := msgFile.Close()
		if err != nil {
			crit("Error closing info log: %s", err.Error())
		}
	}

	if errorFile != nil {
		slog.SetErrorFile(os.Stderr)
		err := errorFile.Close()
		if err != nil {
			crit("Error closing error log: %s", err.Error())
		}
	}
}

func info(f string, v ...interface{}) {
	slog.TMsg(f, v...)
}

func crit(f string, v ...interface{}) {
	slog.TError(f, v...)
}
