package main

import (
	"os"
	"log"
	"github.com/virtual-designer/urnetwork-client-gnu/widgets"
	"github.com/virtual-designer/urnetwork-client-gnu/core"
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

func onActivate(app *gtk.Application) {
	authManager, err := core.NewAuthManager("")

	if err == nil {
		log.Printf("JWT: %s", authManager.Jwt)
	}

	window := widgets.NewAppWindow(authManager)
	app.AddWindow(window.Cast().(*gtk.Window))
	window.Show()
}
