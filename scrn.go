package main

import (
	"bytes"
	"crypto/sha256"
	"errors"
	"fmt"
	"github.com/gotk3/gotk3/gdk"
	"github.com/gotk3/gotk3/gtk"
	"github.com/joho/godotenv"
	"io/ioutil"
	"log"
	"mime/multipart"
	"net/http"
	"os"
)

const (
	HeavyTxtTplSize = 80000
)

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
}

func setupWindow(title string) *gtk.Window {
	win, err := gtk.WindowNew(gtk.WINDOW_TOPLEVEL)
	if err != nil {
		log.Fatal("Unable to create window:", err)
	}
	win.SetTitle(title)
	if _, err := win.Connect("destroy", func() { gtk.MainQuit() }); err != nil {
		_, _ = fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}
	win.SetDefaultSize(500, 300)
	win.SetPosition(gtk.WIN_POS_CENTER)
	return win
}

func setupBox(orient gtk.Orientation) *gtk.Box {
	box, err := gtk.BoxNew(orient, 0)
	if err != nil {
		log.Fatal("Unable to create box:", err)
	}
	return box
}

func setupTview() *gtk.TextView {
	tv, err := gtk.TextViewNew()
	if err != nil {
		log.Fatal("Unable to create TextView:", err)
	}
	return tv
}

func setupBtn(label string, onClick func()) *gtk.Button {
	btn, err := gtk.ButtonNewWithLabel(label)
	if err != nil {
		log.Fatal("Unable to create button:", err)
	}
	_, _ = btn.Connect("clicked", onClick)
	return btn
}

func getClipboard() *gtk.Clipboard {
	display, err := gdk.DisplayGetDefault()
	if err != nil {
		log.Fatal(err)
	}
	clipboard, err := gtk.ClipboardGetForDisplay(display, gdk.SELECTION_CLIPBOARD)
	if err != nil {
		log.Fatal(err)
	}
	return clipboard
}

func main() {
	gtk.Init(nil)

	win := setupWindow("Homemade Screenshotter")
	box := setupBox(gtk.ORIENTATION_VERTICAL)
	win.Add(box)

	tv := setupTview()
	setTextInTview(tv, "Copy text or image and press \"Upload\"")
	box.PackStart(tv, true, true, 0)

	btn := setupBtn("Send To Server", func() {
		url := os.Getenv("UPLOAD_URL")

		clipboard := getClipboard()
		if textContent, textErr := clipboard.WaitForText(); textErr == nil {
			if fileUrl, err := sendTextToServer(textContent, url); err != nil {
				setTextInTview(tv, fmt.Sprintf("Error uploading TEXT.\nDetails: \"%s\"", err.Error()))
				return
			} else {
				setTextInTview(tv, "TXT> "+fileUrl)
				clipboard.SetText(fileUrl)
			}
		} else {
			if imageContent, imageErr := clipboard.WaitForImage(); imageErr == nil {
				if fileUrl, err := sendImageToServer(imageContent, url); err != nil {
					setTextInTview(tv, fmt.Sprintf("Error uploading IMAGE.\nDetails: \"%s\"", err.Error()))
					return
				} else {
					setTextInTview(tv, "IMG> "+fileUrl)
					clipboard.SetText(fileUrl)
				}
			}
		}
	})
	box.Add(btn)

	// Recursively show all widgets contained in this window.
	win.ShowAll()

	// Begin executing the GTK main loop.  This blocks until gtk.MainQuit() is run.
	gtk.Main()
}

func sendTextToServer(text string, url string) (string, error) {
	tplFilename := "./template_heavy.html"
	if len(text) < HeavyTxtTplSize {
		tplFilename = "./template_light.html"
	}
	tplText, err := ioutil.ReadFile(tplFilename)
	if err != nil {
		return "", err
	}
	tplText = bytes.ReplaceAll(tplText, []byte("#CONTENT#"), []byte(text))

	if fileUrl, err := postFile(tplText, "html", url); err != nil {
		return "", err
	} else {
		return fileUrl, nil
	}
}

func sendImageToServer(im *gdk.Pixbuf, url string) (string, error) {
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
	hash := sha256.Sum256(append(content, []byte(os.Getenv("ACCESS_KEY"))...))
	count, err = fw.Write([]byte(fmt.Sprintf("%x", hash)))
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

func getBufferFromTview(tv *gtk.TextView) *gtk.TextBuffer {
	buffer, err := tv.GetBuffer()
	if err != nil {
		log.Fatal("Unable to get buffer:", err)
	}
	return buffer
}

func setTextInTview(tv *gtk.TextView, text string) {
	buffer := getBufferFromTview(tv)
	buffer.SetText(text)
}

func getTextFromTview(tv *gtk.TextView) string {
	buffer := getBufferFromTview(tv)
	start, end := buffer.GetBounds()

	text, err := buffer.GetText(start, end, true)
	if err != nil {
		log.Fatal("Unable to get text:", err)
	}
	return text
}
