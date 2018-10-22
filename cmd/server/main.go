package main

import (
	"context"
	"flag"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
	"github.com/romanyx/integral_db/internal/get"
	"github.com/romanyx/integral_db/internal/set"
	"github.com/romanyx/integral_db/internal/storage"
)

const (
	shutdownTimeout = 30 * time.Second
	readTimeout     = 15 * time.Second
	writeTimeout    = 15 * time.Second
)

func main() {
	var (
		httpAddr    = flag.String("http", "0.0.0.0:80", "HTTP service address.")
		keyLiveTime = flag.Duration("key-live-time", time.Second*30, "key liveness time")
	)

	flag.Parse()

	errChan := make(chan error)

	httpServer := http.Server{
		ReadTimeout:    readTimeout,
		WriteTimeout:   writeTimeout,
		MaxHeaderBytes: 1 << 20,
		Handler:        httpMux(storage.New(), *keyLiveTime),
		Addr:           *httpAddr,
	}

	go func() {
		errChan <- httpServer.ListenAndServe()
	}()

	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)

	for {
		select {
		case err := <-errChan:
			if err != nil {
				log.Fatalf("fatal error: %s", err)
			}
		case s := <-signalChan:
			log.Printf("captured %v. exiting...", s)

			ctx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
			defer cancel()

			if err := httpServer.Shutdown(ctx); err != nil {
				log.Printf("graceful shutdown did not complete in %v : %v", shutdownTimeout, err)
				if err := httpServer.Close(); err != nil {
					log.Fatalf("could not stop http server: %v", err)
				}
			}
		}
	}
}

func httpMux(s storage.Storage, keyLiveTime time.Duration) http.Handler {
	mux := mux.NewRouter()

	postSet := set.NewHandler(set.NewService(s, keyLiveTime))
	mux.HandleFunc("/set", postSet).Methods("POST")
	getGet := get.NewHandler(get.NewService(s))
	mux.HandleFunc("/get", getGet).Methods("GET")

	return mux
}
