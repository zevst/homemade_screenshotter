package main

import (
	"context"
	"crypto/sha256"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
)

var (
	webServer   *http.Server
	letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ_0123456789")
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
	router := mux.NewRouter()

	router.Use(onPanic)

	router.HandleFunc("/", defaultAction)
	router.HandleFunc("/upload", uploadAction)

	webServer = &http.Server{Addr: os.Getenv("LISTEN_ADDR"), Handler: router}

	go func() {
		if err := webServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			panic(err)
		}
	}()

	log.Printf("Run web server: %v", os.Getenv("LISTEN_ADDR"))
}

// onPanic is Mux's middleware to handle panic situations
func onPanic(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer onActionPanic(w)

		next.ServeHTTP(w, r)
	})
}

// panic handler for controller actions
func onActionPanic(w http.ResponseWriter) {
	if r := recover(); r != nil {
		var message string
		switch x := r.(type) {
		case string:
			message = x
		case error:
			message = x.Error()
		default:
			message = fmt.Sprintf("unknown error '%s'", x)
		}

		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte(message))
		fmt.Fprint(os.Stderr, message)
	}
}

// defaultAction shows any std response as "index" action
func defaultAction(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("noaction"))
}

// uploadAction handles uploading, saving a file and returning back URL to file
func uploadAction(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		defaultAction(w, r)
		return
	}
	if err := r.ParseMultipartForm(10 << 20); err != nil {
		log.Panic(err)
	}

	content := r.FormValue("content")
	if content == "" {
		defaultAction(w, r)
	}
	extension := r.FormValue("extension")
	if extension == "" {
		extension = "html"
	}
	accKey := r.FormValue("access_key")
	hash := sha256.Sum256(append([]byte(content), []byte(os.Getenv("ACCESS_KEY"))...))
	if accKey != fmt.Sprintf("%x", hash) {
		w.WriteHeader(http.StatusForbidden)
		_, _ = w.Write([]byte("Access Denied"))
		return
	}

	filename := fmt.Sprintf(
		"%s_%s.%s",
		time.Now().Format("02.01.2006"),
		randStringRunes(15),
		extension,
	)

	f, err := os.OpenFile(os.Getenv("IMAGE_PATH")+"/"+filename, os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte(err.Error()))
		return
	}
	defer f.Close()
	_, err = f.Write([]byte(content))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte(err.Error()))
		return
	}
	w.WriteHeader(http.StatusOK) // this status is expected to be success
	_, _ = w.Write([]byte(os.Getenv("STATIC_SERVER_PATH") + filename))
}

// randStringRunes generates a random string of given length n
func randStringRunes(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
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
