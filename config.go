package main

import (
	"os"

	"github.com/Urethramancer/cross"
)

const (
	// program is used as the base name for some paths.
	program    = "TriggerWarning"
	configname = "config.json"
)

// Version is injected via build flags from git tags.
var Version = "undefined"

// Config is used to parse the main configuration file.
type Config struct {
	// APIPath is the path before each endpoint.
	APIPath string `json:"apipath"`
	// CleanupInterval is how often to clear out stale authentication data.
	CleanupInterval string `json:"cleanupinterval"`
	// Users is the path where users are stored.
	Users string `json:"users"`
	// Triggers is the path where triggers (actions to notify about) are stored.
	Triggers string `json:"triggers"`
	// Watchers is the path where watchers (users to be notified) are stored.
	Watchers string `json:"watchers"`
	// Repos is the path where repositories are cloned.
	Repos string `json:"repos"`
	// Logs settings.
	Logs `json:"logs"`
	// Security settings.
	Security `json:"security"`
	// SSH and git settings.
	SSH `json:"ssh"`
	// Web settings.
	Web `json:"web"`
}

// SSH options for authentication.
type SSH struct {
	// Username for git access (default: "git").
	Username string `json:"username"`
	// PrivateKey to get public key from.
	PrivateKey string `json:"privatekey"`
}

// Logs hold the names of log files to output to.
type Logs struct {
	// Info is for general messages.
	Info string `json:"info,omitempty"`
	// Error is for serious warnings.
	Error string `json:"error,omitempty"`
}

// Security defines password security level and many SSL settings.
type Security struct {
	// Certificate is the cert.pem file.
	Certificate string `json:"certificate,omitempty"`
	// Key is the key.pem counterpart to the cert.
	Key string `json:"key,omitempty"`
	// SSL enables the secure web server instead. A certificate pair is required for it to work.
	SSL bool `json:"ssl,omitempty"`
}

// Web defines addresses, ports and domains to route.
type Web struct {
	// Address is the IP address to bind to. Valid entries are any reachable address,
	// or 0.0.0.0 to bind to all public addresses, or even 127.0.0.1 if you rely on
	// a proxy server to expose it to the world.
	Address string `json:"address"`
	// Port to bind HTTP to.
	Port string `json:"port,omitempty"`
	// Domain is the FQDN for the server, and is mostly used for display purposes.
	Domain string `json:"domain"`
	url    string
}

var cfg Config

func defaultConfig() Config {
	return Config{
		APIPath:         "/",
		CleanupInterval: "10m",
		Users:           "users",
		Triggers:        "triggers",
		Watchers:        "watchers",
		Repos:           "repos",
		Logs: Logs{
			Info:  "info.log",
			Error: "error.log",
		},
		Security: Security{
			Certificate: "./cert.pem",
			Key:         "./key.pem",
		},
		Web: Web{
			Address: "127.0.0.1",
			Port:    "11000",
			Domain:  "localhost",
		},
		SSH: SSH{
			Username:   "git",
			PrivateKey: "id_rsa",
		},
	}
}

func loadConfig() {
	var err error

	name := opt.Config
	if name == "" {
		name, err = cross.GetConfigName(program, configname)
		if err != nil {
			crit("Error getting configuration: %s", err.Error())
			os.Exit(2)
		}
	}

	if !cross.Exists(name) {
		cfg = defaultConfig()
		err = SaveJSON(name, cfg)
		if err != nil {
			crit("Error saving default configuration %s: %s", name, err.Error())
			os.Exit(2)
		}
		info("Created default configuration.")
		return
	}

	err = LoadJSON(name, &cfg)
	if err != nil {
		crit("Error loading %s: %s", name, err.Error())
	}

	info("Loaded config from %s", name)

	// The data paths are very important, so bail out entirely if any of them can't be created.
	if cfg.Users != "" && !cross.Exists(cfg.Users) {
		err := os.MkdirAll(cfg.Users, 0700)
		check(err)
	}

	if cfg.Triggers != "" && !cross.Exists(cfg.Triggers) {
		err := os.MkdirAll(cfg.Triggers, 0700)
		check(err)
	}

	if cfg.Watchers != "" && !cross.Exists(cfg.Watchers) {
		err := os.MkdirAll(cfg.Watchers, 0700)
		check(err)
	}

	if cfg.Repos != "" && !cross.Exists(cfg.Repos) {
		err := os.MkdirAll(cfg.Repos, 0700)
		check(err)
	}

	if cfg.SSH.Username == "" {
		cfg.SSH.Username = "git"
	}
}
