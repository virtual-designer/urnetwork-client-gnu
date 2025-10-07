package widgets

import (
	_ "embed"
	"github.com/virtual-designer/urnetwork-client-gnu/core"
	"github.com/diamondburned/gotk4/pkg/gtk/v4"
	"github.com/diamondburned/gotk4/pkg/gdk/v4"
)

//go:embed window.ui
var appWindowUiXML string

//go:embed window.css
var appWindowCSS string

type AppWindow struct {
	*gtk.Window
	authManager *core.AuthManager
	stack *gtk.Stack
	connectedView *ConnectedView
}

func (window *AppWindow) OnShutdown() {
	if window.stack.VisibleChildName() == "ConnectedView" {
		window.connectedView.OnDisconnect()
	}
}

func NewAppWindow(authManager *core.AuthManager) *AppWindow {
	builder := gtk.NewBuilderFromString(appWindowUiXML)

	cssProvider := gtk.NewCSSProvider()
	cssProvider.LoadFromData(appWindowCSS)

	gtkWindow := builder.GetObject("AppWindow").Cast().(*gtk.Window)
	stack := builder.GetObject("ViewStack").Cast().(*gtk.Stack)
	loginView := NewLoginView(authManager, stack)
	connectedView := NewConnectedView(authManager, stack)
	mainView := NewMainView(authManager, stack, connectedView)

	window := & AppWindow {
		Window: gtkWindow,
		authManager: authManager,
		stack: stack,
		connectedView: connectedView,
	}

	gtk.StyleContextAddProviderForDisplay(
		gdk.DisplayGetDefault(),
		cssProvider,
		gtk.STYLE_PROVIDER_PRIORITY_APPLICATION,
	)

	stack.AddNamed(loginView.Cast().(gtk.Widgetter), "LoginView")
	stack.AddNamed(mainView.Cast().(gtk.Widgetter), "MainView")
	stack.AddNamed(connectedView.Cast().(gtk.Widgetter), "ConnectedView")

	if authManager.Jwt == "" {
		stack.SetVisibleChildName("LoginView")
	} else {
		stack.SetVisibleChildName("MainView")
	}

	return window
}
