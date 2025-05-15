package main

import (
	_ "embed"
	"fmt"
	"os"

	"github.com/gotk3/gotk3/gtk"
)

//go:embed style.css
var styleCss string

func main() {
	gtk.Init(&os.Args)
	window, _ := gtk.WindowNew(gtk.WINDOW_TOPLEVEL)
	window.SetTitle("Gotk3 Example")
	window.Connect("destroy", func() {
		gtk.MainQuit()
	})

	cssProvider, _ := gtk.CssProviderNew()
	cssProvider.LoadFromData(styleCss)
	gtk.AddProviderForScreen(window.GetScreen(), cssProvider, gtk.STYLE_PROVIDER_PRIORITY_USER)

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
