package widgets

import (
	_ "embed"
	"github.com/diamondburned/gotk4/pkg/gtk/v4"
	"github.com/diamondburned/gotk4/pkg/gdk/v4"
)

//go:embed main.ui
var mainUiXML string

//go:embed main.css
var mainCSS string

type MainView struct {
	*gtk.Box
}

func NewMainView() *MainView {
	builder := gtk.NewBuilderFromString(mainUiXML)

	cssProvider := gtk.NewCSSProvider()
	cssProvider.LoadFromData(mainCSS)

	view := builder.GetObject("MainView").Cast().(*gtk.Box)

	mainView := & MainView {
		Box: view,
	}

	gtk.StyleContextAddProviderForDisplay(
		gdk.DisplayGetDefault(),
		cssProvider,
		gtk.STYLE_PROVIDER_PRIORITY_APPLICATION,
	)

	return mainView
}
