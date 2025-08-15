package main

import (
	"os"
	"github.com/virtual-designer/urnetwork-client-gnu/widgets"
	"github.com/diamondburned/gotk4/pkg/gtk/v4"
	"github.com/diamondburned/gotk4/pkg/gio/v2"
)

const APP_ID = "me.rakinar2.gnu.linux.urnetwork.client"

func main() {
	app := gtk.NewApplication(APP_ID, gio.ApplicationFlagsNone)

	app.ConnectActivate(func () {
		onActivate(app)
	})

	if code := app.Run(os.Args); code > 0 {
		os.Exit(code)
	}
}

func onActivate (app *gtk.Application) {
	window := widgets.NewLoginWindow()
	app.AddWindow(window.Cast().(*gtk.Window))
	window.Show()
}
