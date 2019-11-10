package actions

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"os"
	"time"
)

const (
	randStringSlugLength = 15
)

var (
	randStringSlugRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ_0123456789")
)

// Upload handles uploading, saving a file and returning back URL to file
func Upload(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		Default(w, r)
		return
	}
	if err := r.ParseMultipartForm(10 << 20); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_, _ = w.Write([]byte("Bad Request"))
		return
	}

	content := r.FormValue("content")
	if content == "" {
		Default(w, r)
	}
	extension := r.FormValue("extension")
	if extension == "" {
		extension = "html"
	}
	accKey := r.FormValue("access_key")
	if accKey != hashPayload([]byte(content)) {
		w.WriteHeader(http.StatusForbidden)
		_, _ = w.Write([]byte("Access Denied"))
		return
	}

	filename := fmt.Sprintf(
		"%s_%s.%s",
		time.Now().Format("02.01.2006"),
		randStringRunes(randStringSlugLength),
		extension,
	)

	f, err := os.OpenFile(os.Getenv("IMAGE_PATH")+"/"+filename, os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		log.Panic(err)
	}
	defer func(f *os.File) { _ = f.Close() }(f)
	_, err = f.Write([]byte(content))
	if err != nil {
		log.Panic(err)
	}
	w.WriteHeader(http.StatusOK) // this status is expected to be success
	_, _ = w.Write([]byte(os.Getenv("STATIC_SERVER_PATH") + filename))
}

func hashPayload(payload []byte) string {
	hasher := hmac.New(sha256.New, []byte(os.Getenv("ACCESS_KEY")))
	hasher.Write(payload)
	hash := hasher.Sum(nil)
	return hex.EncodeToString(hash)
}

// randStringRunes generates a random string of given length n
func randStringRunes(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = randStringSlugRunes[rand.Intn(len(randStringSlugRunes))]
	}
	return string(b)
}
