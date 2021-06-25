package main

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
)

const (
	FATAL = "FATAL: "
	ERROR = "ERROR: "
	WARN  = "WARN: "
	INFO  = "INFO: "
	DEBUG = "DEBUG: "
)

func handlerFunc(w http.ResponseWriter, r *http.Request) {
	addr := strings.Split(r.Header.Get("X-Forwarded-For"), ",")[0]
	if addr == "" {
		host, _, err := net.SplitHostPort(r.RemoteAddr)
		if err != nil {
			errorResponse(w, http.StatusOK, err)
			return

		}
		addr = host
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, `{"status": %d, "address": %q}`, http.StatusOK, addr)
}

func errorResponse(w http.ResponseWriter, status int, err error) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(status)
	fmt.Fprintf(w, `{"status": %d, "message": %q}`, status, err.Error())
}

func main() {
	log.SetFlags(log.Ldate | log.Ltime | log.Lmicroseconds | log.Lshortfile)
	sig := make(chan os.Signal)
	signal.Notify(sig, syscall.SIGINT)
	signal.Notify(sig, syscall.SIGTERM)
	signal.Notify(sig, syscall.SIGKILL)
	go func() {
		log.Print(INFO, <-sig, ": received signal")
		os.Exit(0)
	}()

	http.HandleFunc("/api/global-ip", handlerFunc)

	log.Print(INFO, "server start")
	err := http.ListenAndServe(":8008", nil)
	if err != nil {
		log.Print(FATAL, "http.ListenAndServe(): ", err)
		os.Exit(1)
	}

	/*
		crt := filepath.Join(os.Getenv("GOPATH"), "crt", "server.crt")
		key := filepath.Join(os.Getenv("GOPATH"), "crt", "server.key")

		err := http.ListenAndServeTLS(":8443", crt, key, nil)
		if err != nil {
			log.Print(FATAL, "ListenAndServeTLS: ", err)
			os.Exit(1)
		}
	*/
}
