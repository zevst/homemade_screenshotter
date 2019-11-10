package ui

import (
	"fmt"
	"github.com/gotk3/gotk3/gdk"
	"github.com/gotk3/gotk3/gtk"
	"log"
	"os"
)

func SetupWindow(title string) *gtk.Window {
	win, err := gtk.WindowNew(gtk.WINDOW_TOPLEVEL)
	if err != nil {
		log.Fatal("Unable to create window:", err)
	}
	win.SetTitle(title)
	if _, err := win.Connect("destroy", func() { gtk.MainQuit() }); err != nil {
		_, _ = fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}
	wd, err := os.Getwd()
	if err != nil {
		_, _ = fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}
	fmt.Println(wd)
	if err := win.SetIconFromFile(wd + "/icon.png"); err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "Error getting icon %q \r\n", err.Error())
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

func SetupTview() *gtk.TextView {
	tv, err := gtk.TextViewNew()
	if err != nil {
		log.Fatal("Unable to create TextView:", err)
	}
	return tv
}

func SetupBtn(label string, onClick func()) *gtk.Button {
	btn, err := gtk.ButtonNewWithLabel(label)
	if err != nil {
		log.Fatal("Unable to create button:", err)
	}
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
