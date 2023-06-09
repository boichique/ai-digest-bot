package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"digest_bot_database/internal/config"
	"digest_bot_database/internal/crontasks"
	"digest_bot_database/internal/server"

	"github.com/robfig/cron"
)

const (
	dbConnectTimeout = 10 * time.Second
	gracefulTimeout  = 10 * time.Second
)

func main() {
	cfg, err := config.NewConfig()
	failOnError(err, "parse config")

	srv, err := server.New(context.Background(), cfg)
	failOnError(err, "create server")

	cr := cron.New()
	cr.AddFunc("* 10 * * *", func() { crontasks.UpdateFullDigestsForUsers() }) // update every hour on 10 minutes (* 10 * * *)
	cr.Start()

	go func() {
		sigCh := make(chan os.Signal, 1)
		signal.Notify(sigCh, syscall.SIGINT)
		<-sigCh
		log.Println("received interrupt signal. Shutting down...")

		ctx, cancel := context.WithTimeout(context.Background(), gracefulTimeout)
		defer cancel()

		if err := srv.Shutdown(ctx); err != nil {
			log.Fatalf("shutdown server: %s", err)
		}
	}()

	if err = srv.Start(); err != http.ErrServerClosed {
		log.Fatal(err)
	}
}

func failOnError(err error, message string) {
	if err != nil {
		log.Fatalf("%s: %s", message, err)
	}
}
