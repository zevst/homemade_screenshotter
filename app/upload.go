package app

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"github.com/gotk3/gotk3/gdk"
	"html"
	"io/ioutil"
	"log"
	"mime/multipart"
	"net/http"
	"os"
)

const (
	HeavyTxtTplSize = 80000
)

func postFile(content []byte, extension string, targetURL string) (string, error) {
	var b bytes.Buffer
	w := multipart.NewWriter(&b)
	fw, err := w.CreateFormField("content")
	if err != nil {
		return "", err
	}
	count, err := fw.Write(content)
	if err != nil {
		return "", err
	}
	if count == 0 {
		return "", errors.New("nil-size file")
	}
	fw, err = w.CreateFormField("extension")
	if err != nil {
		log.Fatal(err)
	}
	count, err = fw.Write([]byte(extension))
	if err != nil {
		return "", err
	}
	if count == 0 {
		return "", errors.New("empty extension specified")
	}
	// security hdr add
	fw, err = w.CreateFormField("access_key")
	if err != nil {
		return "", err
	}
	count, err = fw.Write([]byte(hashPayload(content)))
	if err != nil {
		return "", err
	}
	if count == 0 {
		return "", errors.New("security token not specified")
	}

	if err := w.Close(); err != nil {
		return "", err
	}

	req, err := http.NewRequest(http.MethodPost, targetURL, &b)
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", w.FormDataContentType())
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}
	body, _ := ioutil.ReadAll(res.Body)
	if res.StatusCode != http.StatusOK { // something wrong
		return "", errors.New(fmt.Sprintf("Error while uploading file:\nStatus: %s\nMessage: %s", res.Status, string(body)))
	}
	return string(body), nil
}

// hashPayload calculates HMAC sha256 of given data; key is in .env file
func hashPayload(payload []byte) string {
	hasher := hmac.New(sha256.New, []byte(os.Getenv("ACCESS_KEY")))
	hasher.Write(payload)
	hash := hasher.Sum(nil)
	return hex.EncodeToString(hash)
}

// SendTextToServer sends data to server as text
func SendTextToServer(text string, url string) (string, error) {
	tplFilename := "./templates/heavy.html"
	if len(text) < HeavyTxtTplSize {
		tplFilename = "./templates/light.html"
	}
	tplText, err := ioutil.ReadFile(tplFilename)
	if err != nil {
		return "", err
	}

	text = html.EscapeString(text)
	tplText = bytes.ReplaceAll(tplText, []byte("#CONTENT#"), []byte(text))

	if fileUrl, err := postFile(tplText, "html", url); err != nil {
		return "", err
	} else {
		return fileUrl, nil
	}
}

// SendImageToServer sends data to server as PNG image
func SendImageToServer(im *gdk.Pixbuf, url string) (string, error) {
	filename := os.Getenv("TMP_FOLDER") + "/temp_img_buf"
	if err := im.SavePNG(filename, 9); err != nil {
		return "", err
	}
	content, err := ioutil.ReadFile(filename)
	if err != nil {
		return "", err
	}
	if fileUrl, err := postFile(content, "png", url); err != nil {
		return "", err
	} else {
		return fileUrl, nil
	}
}
