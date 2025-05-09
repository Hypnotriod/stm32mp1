package main

import (
	"fmt"

	"github.com/gotk3/gotk3/gtk"
)

func main() {
	gtk.Init(nil)
	window, _ := gtk.WindowNew(gtk.WINDOW_TOPLEVEL)
	window.SetTitle("Gotk3 Example")
	window.Connect("destroy", func() {
		gtk.MainQuit()
	})

	buttonBox, _ := gtk.ButtonBoxNew(gtk.ORIENTATION_HORIZONTAL)

	button, _ := gtk.ButtonNewWithLabel("Hello, gotk3!")
	button.Connect("clicked", func() {
		fmt.Println("Hello, gotk3!")
		gtk.MainQuit()
	})

	buttonBox.Add(button)

	window.Add(buttonBox)
	window.SetDefaultSize(300, 200)
	window.ShowAll()

	gtk.Main()
}
