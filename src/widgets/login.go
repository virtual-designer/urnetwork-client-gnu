package widgets

import (
	_ "embed"
	"log"
	"github.com/virtual-designer/urnetwork-client-gnu/core"
	"github.com/diamondburned/gotk4/pkg/gtk/v4"
	"github.com/diamondburned/gotk4/pkg/gdk/v4"
)

//go:embed login.ui
var loginUiXML string

//go:embed login.css
var loginCSS string

type LoginView struct {
	*gtk.Box
	stack *gtk.Stack
	authManager *core.AuthManager
	loginButton *gtk.Button
	emailEntry *gtk.Entry
	passwordEntry *gtk.PasswordEntry
	errorLabel *gtk.Label
}

func NewLoginView(authManager *core.AuthManager, stack *gtk.Stack) *LoginView {
	builder := gtk.NewBuilderFromString(loginUiXML)

	cssProvider := gtk.NewCSSProvider()
	cssProvider.LoadFromData(loginCSS)

	view := builder.GetObject("LoginView").Cast().(*gtk.Box)
	loginButton := builder.GetObject("LoginButton").Cast().(*gtk.Button)
	emailEntry := builder.GetObject("EmailEntry").Cast().(*gtk.Entry)
	passwordEntry := builder.GetObject("PasswordEntry").Cast().(*gtk.PasswordEntry)
	errorLabel := builder.GetObject("ErrorLabel").Cast().(*gtk.Label)

	loginView := & LoginView {
		Box: view,
		stack: stack,
		authManager: authManager,
		loginButton: loginButton,
		emailEntry: emailEntry,
		passwordEntry: passwordEntry,
		errorLabel: errorLabel,
	}

	gtk.StyleContextAddProviderForDisplay(
		gdk.DisplayGetDefault(),
		cssProvider,
		gtk.STYLE_PROVIDER_PRIORITY_APPLICATION,
	)

	loginButton.Connect("clicked", func() {
		email := emailEntry.Text()
		password := passwordEntry.Text()

		log.Println("Attempting Login")

		result, err := authManager.PerformAuth(email, password)

		if err != nil {
			log.Printf("Login error: %s", err.Error())
			msg := err.Error()

			if msg == "" {
				msg = "An error has occurred"
			}

			errorLabel.SetLabel(msg)
			errorLabel.SetVisible(true)
		} else {
			errorLabel.SetVisible(false)
			log.Printf("Result: %#v\n", result)
			stack.SetVisibleChildName("MainView")
		}
	})

	return loginView
}
