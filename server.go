package main

import (
	"context"
	"fmt"
	"homemadeScreenshotterUploader/controller"
	"log"
	"math/rand"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/joho/godotenv"
)

var (
	webServer *http.Server
)

func init() {
	rand.Seed(time.Now().UnixNano())

	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
}

func main() {
	stopch := make(chan os.Signal, 1)
	signal.Notify(stopch, syscall.SIGINT, syscall.SIGTERM)

	handle()

	<-stopch
	shutdown()
}

// handle - main HTTP handler
func handle() {
	webServer = &http.Server{Addr: os.Getenv("LISTEN_ADDR"), Handler: controller.Router()}

	go func() {
		if err := webServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			panic(err)
		}
	}()

	log.Printf("Run web server: %v", os.Getenv("LISTEN_ADDR"))
}

// shutdown function performs graceful server shutdown
func shutdown() {
	log.Println("Service shutdown initiated")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := webServer.Shutdown(ctx); err != nil {
		_, _ = fmt.Fprint(os.Stderr, err)
	} else {
		log.Println("Web server successfully stopped")
	}
}
