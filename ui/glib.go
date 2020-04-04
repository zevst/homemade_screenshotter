package ui

import (
	"github.com/gotk3/gotk3/glib"
	"github.com/gotk3/gotk3/gtk"
	"log"
)

// SetTextGlib sets text to label
// uses glib.IdleAdd (allowing to run code in goroutine)
func SetTextGlib(label *gtk.Label, text string) {
	if _, err := glib.IdleAdd(label.SetText, text); err != nil {
		log.Fatalf("SetTextGlib failed: %q", err.Error())
	}
}

// PrependMarkupGlib prepends label text with markup
// uses glib.IdleAdd (allowing to run code in goroutine)
func PrependMarkupGlib(label *gtk.Label, text string) {
	oldText, err := label.GetText()
	if err != nil {
		log.Fatalf("Unable to get text: %q", err.Error())
	}
	newText := text + "\r\n" + oldText
	if _, err := glib.IdleAdd(label.SetMarkup, newText); err != nil {
		log.Fatalf("PrependMarkupGlib failed: %q", err.Error())
	}
}

// SetClipboardTextGlib sets text to clipboard
// uses glib.IdleAdd (allowing to run code in goroutine)
func SetClipboardTextGlib(clipboard *gtk.Clipboard, text string) {
	if _, err := glib.IdleAdd(clipboard.SetText, text); err != nil {
		log.Fatalf("SetClipboardTextGlib failed: %q", err.Error())
	}
}
