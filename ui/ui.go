package ui

import (
	"github.com/gotk3/gotk3/gdk"
	"github.com/gotk3/gotk3/gtk"
	"log"
)

func SetupWindow(application *gtk.Application, title string, mainFolder string) *gtk.ApplicationWindow {
	win, err := gtk.ApplicationWindowNew(application)
	if err != nil {
		log.Fatal("Unable to create window:", err)
	}
	win.SetTitle(title)
	if _, err := win.Connect("destroy", func() { gtk.MainQuit() }); err != nil {
		log.Fatal(err)
	}
	if err := win.SetIconFromFile(mainFolder + "/icon.png"); err != nil {
		log.Fatal(err)
	}
	win.SetDefaultSize(700, 500)
	win.SetResizable(false)
	win.SetPosition(gtk.WIN_POS_CENTER)
	return win
}

func SetupBox(orient gtk.Orientation) *gtk.Box {
	box, err := gtk.BoxNew(orient, 0)
	if err != nil {
		log.Fatal("Unable to create box:", err)
	}
	return box
}

// SetupTview creates TextView widget and a scrollable wrapper for it and return both of them. Crashes an app if any error
func SetupTview() (*gtk.TextView, *gtk.ScrolledWindow) {
	// create wrapper
	scw, err := gtk.ScrolledWindowNew(nil, nil)
	if err != nil {
		log.Fatal("Unable to create TextView scrollable wrapper: ", err)
	}
	scw.SetSizeRequest(700, 230)

	// create TextView widget
	tv, err := gtk.TextViewNew()
	if err != nil {
		log.Fatal("Unable to create TextView: ", err)
	}
	//tv.SetSizeRequest(700, 230)
	tv.SetEditable(false)
	//tv.SetWrapMode(gtk.WRAP_WORD)
	tv.SetMarginBottom(5)

	scw.Add(tv)
	return tv, scw
}

// SetupBtn creates Button to start uploading on server ...
func SetupBtn(label string, onClick func()) *gtk.Button {
	btn, err := gtk.ButtonNewWithLabel(label)
	if err != nil {
		log.Fatal("Unable to create button:", err)
	}
	btn.SetSizeRequest(700, 40)
	_, _ = btn.Connect("clicked", onClick)
	return btn
}

func SetTextInTview(tv *gtk.TextView, text string) {
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

func PrependTextInTview(tv *gtk.TextView, text string) {
	SetTextInTview(tv, text+"\r\n"+getTextFromTview(tv))
}

func getBufferFromTview(tv *gtk.TextView) *gtk.TextBuffer {
	buffer, err := tv.GetBuffer()
	if err != nil {
		log.Fatal("Unable to get buffer:", err)
	}
	return buffer
}

func GetClipboard() *gtk.Clipboard {
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
