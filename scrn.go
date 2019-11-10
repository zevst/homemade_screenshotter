package main

import (
	"fmt"
	"github.com/gotk3/gotk3/gtk"
	"github.com/joho/godotenv"
	"homemadeScreenshotter/app"
	"homemadeScreenshotter/ui"
	"log"
	"os"
)

var (
	InstallFld string
	textView   *gtk.TextView
)

func init() {
	if InstallFld != "" {
		if err := os.Chdir(InstallFld); err != nil {
			log.Fatalf("Can not chdir to %q", InstallFld)
		}
	}
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}
}

func main() {
	gtk.Init(nil)

	win := ui.SetupWindow("Homemade Screenshotter")
	box := ui.SetupBox(gtk.ORIENTATION_VERTICAL)
	win.Add(box)

	textView = ui.SetupTview()
	ui.SetTextInTview(textView, "Copy text or image and press \"Upload\"")
	box.PackStart(textView, true, true, 0)

	btn := ui.SetupBtn("Upload", doUpload)
	box.Add(btn)

	// Recursively show all widgets contained in this window.
	win.ShowAll()

	// Begin executing the GTK main loop.  This blocks until gtk.MainQuit() is run.
	gtk.Main()
}

// doUpload performs uploading text/img to server; triggers by click on button
func doUpload() {
	url := os.Getenv("UPLOAD_URL")

	clipboard := ui.GetClipboard()
	if textContent, textErr := clipboard.WaitForText(); textErr == nil {
		if fileUrl, err := app.SendTextToServer(textContent, url); err != nil {
			ui.SetTextInTview(textView, fmt.Sprintf("Error uploading TEXT.\nDetails: \"%s\"", err.Error()))
			return
		} else {
			ui.SetTextInTview(textView, "TXT> "+fileUrl)
			clipboard.SetText(fileUrl)
		}
	} else {
		if imageContent, imageErr := clipboard.WaitForImage(); imageErr == nil {
			if fileUrl, err := app.SendImageToServer(imageContent, url); err != nil {
				ui.SetTextInTview(textView, fmt.Sprintf("Error uploading IMAGE.\nDetails: \"%s\"", err.Error()))
				return
			} else {
				ui.SetTextInTview(textView, "IMG> "+fileUrl)
				clipboard.SetText(fileUrl)
			}
		}
	}
}
