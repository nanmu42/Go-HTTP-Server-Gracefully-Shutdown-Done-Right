package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	downright "github.com/nanmu42/Go-HTTP-Server-Gracefully-Shutdown-Done-Right"
)

var (
	port           = flag.Int("port", 3100, "port to listen on")
	sleepSeconds   = flag.Int("sleep", 6, "seconds to sleep before response")
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

	server := &GracefulServer{
		Server: &http.Server{
			Addr:    fmt.Sprintf(":%d", *port),
			Handler: downright.SlowHandler(*sleepSeconds),
		},
	}

	go server.WaitForExitingSignal(time.Duration(*timeoutSeconds) * time.Second)

	log.Printf("listening on port %d...", *port)
	err = server.ListenAndServe()
	if err != nil {
		err = fmt.Errorf("unexpected error from ListenAndServe: %w", err)
	}
	log.Println("main goroutine exited.")
}

type GracefulServer struct {
	Server           *http.Server
	shutdownFinished chan struct{}
}

func (s *GracefulServer) ListenAndServe() (err error) {
	if s.shutdownFinished == nil {
		s.shutdownFinished = make(chan struct{})
	}

	err = s.Server.ListenAndServe()
	if err == http.ErrServerClosed {
		// expected error after calling Server.Shutdown().
		err = nil
	} else if err != nil {
		err = fmt.Errorf("unexpected error from ListenAndServe: %w", err)
		return
	}

	log.Println("waiting for shutdown finishing...")
	<-s.shutdownFinished
	log.Println("shutdown finished")

	return
}

func (s *GracefulServer) WaitForExitingSignal(timeout time.Duration) {
	var waiter = make(chan os.Signal, 1) // buffered channel
	signal.Notify(waiter, syscall.SIGTERM, syscall.SIGINT)

	// blocks here until there's a signal
	<-waiter

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	err := s.Server.Shutdown(ctx)
	if err != nil {
		log.Println("shutting down: " + err.Error())
	} else {
		log.Println("shutdown processed successfully")
		close(s.shutdownFinished)
	}
}
