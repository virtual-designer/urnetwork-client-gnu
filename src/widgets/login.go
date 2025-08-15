package widgets

import (
	_ "embed"
	"github.com/diamondburned/gotk4/pkg/gtk/v4"
	"github.com/diamondburned/gotk4/pkg/gdk/v4"
)

//go:embed login.ui
var loginUiXML string

//go:embed login.css
var loginCSS string

type LoginWindow struct {
	*gtk.Window
}

func NewLoginWindow() *LoginWindow {
	builder := gtk.NewBuilderFromString(loginUiXML)

	cssProvider := gtk.NewCSSProvider()
	cssProvider.LoadFromData(loginCSS)

	window := builder.GetObject("LoginWindow").Cast().(*gtk.Window)

	loginWindow := & LoginWindow {
		Window: window,
	}

	gtk.StyleContextAddProviderForDisplay(
		gdk.DisplayGetDefault(),
		cssProvider,
		gtk.STYLE_PROVIDER_PRIORITY_APPLICATION,
	)

	return loginWindow
}
