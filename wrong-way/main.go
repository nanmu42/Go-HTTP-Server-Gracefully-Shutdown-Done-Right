package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	downright "github.com/nanmu42/Go-HTTP-Server-Gracefully-Shutdown-Done-Right"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"
)

var (
	port = flag.Int("port", 3000, "port to listen on")
	sleepSeconds = flag.Int("sleep", 6, "seconds to sleep before response")
	timeoutSeconds = flag.Int("timeout", 10, "seconds to wait before shutting down")
)

func main() {
	flag.Parse()

	var err error
	defer func() {
		if err != nil {
			log.Println("exited with error: " + err.Error())
		}
	}()

	if *sleepSeconds < 0 {
		err = errors.New("you can't reverse the time even if you sleep upside down :P")
		return
	}
	if *timeoutSeconds <= 0 {
		err = fmt.Errorf("timeout must be larger than 0, got %d", *timeoutSeconds)
		return
	}

	server := &http.Server{
		Addr:              fmt.Sprintf(":%d", *port),
		Handler:           downright.SlowHandler(*sleepSeconds),
	}

	go waitForExitingSignal(server, time.Duration(*timeoutSeconds) * time.Second)

	log.Printf("listening on port %d...", *port)
	err = server.ListenAndServe()
	if err == http.ErrServerClosed {
		// expected error after calling Server.Shutdown().
		err = nil
	} else if err != nil {
		err = fmt.Errorf("unexpected error from ListenAndServe: %w", err)
	}
	log.Println("main goroutine exited.")
}

func waitForExitingSignal(server *http.Server, timeout time.Duration) {
	var waiter = make(chan os.Signal, 1) // buffered channel
	signal.Notify(waiter, os.Interrupt)

	// blocks here until there's a signal
	<- waiter

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	err := server.Shutdown(ctx)
	if err != nil {
		log.Println("shutting down: " + err.Error())
	} else {
		log.Println("shutdown processed successfully")
	}
}
