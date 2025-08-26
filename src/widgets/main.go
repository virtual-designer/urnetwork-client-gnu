package widgets

import (
	_ "embed"
	"log"
	"sort"
	"os/exec"
	"github.com/virtual-designer/urnetwork-client-gnu/core"
	"github.com/diamondburned/gotk4/pkg/gtk/v4"
	"github.com/diamondburned/gotk4/pkg/gdk/v4"
)

//go:embed main.ui
var mainUiXML string

//go:embed main.css
var mainCSS string

type MainView struct {
	*gtk.Box
	authManager *core.AuthManager
	loadingBox *gtk.Box
	locationList *gtk.ListBox
	locationListWrapper *gtk.Box
	connectButton *gtk.Button
	errorViewLabel *gtk.Label
	errorViewWrapper *gtk.Box

	connectInfoWrapper *gtk.Box
	connectStatus *gtk.Label
}

func (mainView *MainView) onConnectClick() {
	mainView.connectInfoWrapper.SetVisible(true)
	mainView.connectStatus.SetLabel("Connecting")
}

func (mainView *MainView) loadLocations() {
	errorViewLabel := mainView.errorViewLabel
	errorViewWrapper := mainView.errorViewWrapper
	connectButton := mainView.connectButton
	locationList := mainView.locationList
	locationListWrapper := mainView.locationListWrapper
	loadingBox := mainView.loadingBox

	locations, err := core.GetLocations(mainView.authManager.Jwt)

	if err != nil {
		log.Println("An error has occurred: ", err)
		errorViewLabel.SetLabel(err.Error())

		loadingBox.SetVisible(false)
		locationListWrapper.SetVisible(false)
		errorViewWrapper.SetVisible(true)

		return
	}

	sort.Slice(locations.Locations, func(i, j int) bool {
		return locations.Locations[i].Name < locations.Locations[j].Name
	})

	errorViewWrapper.SetVisible(false)
	loadingBox.SetVisible(false)
	locationListWrapper.SetVisible(true)

	for i, location := range locations.Locations {
		label := gtk.NewLabel(location.Name)
		locationList.Insert(label, i)
	}

	locationList.Connect("row-selected", func() {
		connectButton.SetSensitive(true)
	})

	connectButton.Connect("clicked", func() {
		mainView.onConnectClick()
	})
}

func NewMainView(authManager *core.AuthManager) *MainView {
	builder := gtk.NewBuilderFromString(mainUiXML)

	cssProvider := gtk.NewCSSProvider()
	cssProvider.LoadFromData(mainCSS)

	mainView := & MainView {
		Box: builder.GetObject("MainView").Cast().(*gtk.Box),
		authManager: authManager,
		loadingBox: builder.GetObject("LoadingBox").Cast().(*gtk.Box),
		locationList: builder.GetObject("LocationList").Cast().(*gtk.ListBox),
		locationListWrapper: builder.GetObject("LocationListWrapper").Cast().(*gtk.Box),
		connectButton: builder.GetObject("ConnectButton").Cast().(*gtk.Button),
		errorViewLabel: builder.GetObject("ErrorViewLabel").Cast().(*gtk.Label),
		errorViewWrapper: builder.GetObject("ErrorViewWrapper").Cast().(*gtk.Box),
		connectInfoWrapper: builder.GetObject("ConnectInfoWrapper").Cast().(*gtk.Box),
		connectStatus: builder.GetObject("ConnectStatusLabel").Cast().(*gtk.Label),
	}

	gtk.StyleContextAddProviderForDisplay(
		gdk.DisplayGetDefault(),
		cssProvider,
		gtk.STYLE_PROVIDER_PRIORITY_APPLICATION,
	)

	go mainView.loadLocations()
	return mainView
}
