package main

import (
	"context"
	"crypto/tls"
	"net"
	"net/http"
	"time"

	"github.com/gorilla/mux"
)

// Servers hold the two possible web servers which can be started.
var webserver *http.Server

func initWeb() {
	var address string

	if cfg.Security.SSL {
		address = net.JoinHostPort(cfg.Web.Address, cfg.Web.Port)
		cfg.Web.url = "https://" + cfg.Web.Domain
		if cfg.Web.Port != "443" {
			cfg.Web.url += ":" + cfg.Web.Port
		}
	} else {
		address = net.JoinHostPort(cfg.Web.Address, cfg.Web.Port)
		cfg.Web.url = "http://" + cfg.Web.Domain
		if cfg.Web.Port != "80" {
			cfg.Web.url += ":" + cfg.Web.Port
		}
	}

	router := mux.NewRouter().StrictSlash(true)
	api := "/"
	if cfg.APIPath != "" {
		api += cfg.APIPath + "/"
	}
	sub := router.Host(cfg.Web.Domain).Methods("POST").PathPrefix(api).Subrouter()

	if cfg.Security.SSL {
		secure := sub.Schemes("https").Subrouter()
		initSSLHandlers(secure)
		info("Starting secure web server on %s (%s)", address, cfg.Web.url)
		go startHTTPS(address, secure)
	} else {
		insecure := sub.Schemes("http").Subrouter()
		initHandlers(insecure)
		info("Starting plain web server on %s (%s)", address, cfg.Web.url)
		go startHTTP(address, insecure)
	}
}

func startHTTP(address string, r http.Handler) {
	webserver = &http.Server{
		Addr:    address,
		Handler: r,
		// Safe numbers which should help against DoS.
		IdleTimeout:  time.Second * 30,
		ReadTimeout:  time.Second * 5,
		WriteTimeout: time.Second * 10,
	}
	err := webserver.ListenAndServe()
	if err == http.ErrServerClosed {
		info("HTTP server shut down cleanly.")
	}
}

func startHTTPS(address string, r http.Handler) {
	config := &tls.Config{
		MinVersion: tls.VersionTLS12,
	}
	config.NextProtos = []string{"http/1.1"}
	config.MinVersion = tls.VersionTLS12
	config.CurvePreferences = []tls.CurveID{tls.CurveP521, tls.CurveP384, tls.CurveP256}
	config.PreferServerCipherSuites = true
	config.CipherSuites = []uint16{
		tls.TLS_ECDHE_ECDSA_WITH_AES_256_GCM_SHA384,
		tls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,
		tls.TLS_ECDHE_ECDSA_WITH_CHACHA20_POLY1305,
		tls.TLS_ECDHE_RSA_WITH_CHACHA20_POLY1305,
	}
	webserver = &http.Server{
		Addr:         address,
		Handler:      r,
		IdleTimeout:  time.Second * 30,
		ReadTimeout:  time.Second * 5,
		WriteTimeout: time.Second * 10,
		TLSConfig:    config,
		TLSNextProto: make(map[string]func(*http.Server, *tls.Conn, http.Handler), 0),
	}
	// The certificate should be the fullchain version if using Let's Encrypt.
	err := webserver.ListenAndServeTLS(cfg.Security.Certificate, cfg.Security.Key)
	if err == http.ErrServerClosed {
		info("HTTPS server shut down: %s", err)
	}
}

func stopServers() {
	ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond*500)
	defer cancel()
	err := webserver.Shutdown(ctx)
	if err != nil {
		crit("Error: %s", err.Error())
	}
}

type (
	// Handler is the signature for web route handlers.
	Handler func(w http.ResponseWriter, r *http.Request) error
)

// Handle walks through its chain of middleware for a web request.
func Handle(handlers ...Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Handle handlers until one can't handle it.
		for _, handler := range handlers {
			err := handler(w, r)
			if err != nil {
				crit("Error: %s", err.Error())
				return
			}
		}
	})
}

func addJSONHeaders(w http.ResponseWriter, r *http.Request) error {
	w.Header().Set("Content-Type", "application/json;charset=UTF-8")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	return nil
}

// addSSLHeaders ensures browsers stick to SSL after they start using it.
func addSSLHeaders(w http.ResponseWriter, r *http.Request) error {
	w.Header().Set("Strict-Transport-Security", "max-age=15768000; includeSubDomains")
	w.Header().Set("Content-Type", "charset=UTF-8")
	return nil
}

func initHandlers(r *mux.Router) {
	r.Handle(PathTrigger, Handle(addJSONHeaders, apiTrigger))
	r.Handle(PathWatch, Handle(addJSONHeaders, apiWatch))
	r.Handle(PathUnwatch, Handle(addJSONHeaders, apiUnwatch))
}

func initSSLHandlers(r *mux.Router) {
	r.Handle(PathTrigger, Handle(addSSLHeaders, addJSONHeaders, apiTrigger))
	r.Handle(PathWatch, Handle(addSSLHeaders, addJSONHeaders, apiWatch))
	r.Handle(PathUnwatch, Handle(addSSLHeaders, addJSONHeaders, apiUnwatch))
}
