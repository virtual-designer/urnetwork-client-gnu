package widgets

import (
	_ "embed"
	"log"
	"sort"
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
}

func NewMainView(authManager *core.AuthManager) *MainView {
	builder := gtk.NewBuilderFromString(mainUiXML)

	cssProvider := gtk.NewCSSProvider()
	cssProvider.LoadFromData(mainCSS)

	view := builder.GetObject("MainView").Cast().(*gtk.Box)
	loadingBox := builder.GetObject("LoadingBox").Cast().(*gtk.Box)
	locationList := builder.GetObject("LocationList").Cast().(*gtk.ListBox)
	locationListWrapper := builder.GetObject("LocationListWrapper").Cast().(*gtk.Box)
	connectButton := builder.GetObject("ConnectButton").Cast().(*gtk.Button)

	mainView := & MainView {
		Box: view,
		authManager: authManager,
		loadingBox: loadingBox,
		locationList: locationList,
		locationListWrapper: locationListWrapper,
		connectButton: connectButton,
	}

	gtk.StyleContextAddProviderForDisplay(
		gdk.DisplayGetDefault(),
		cssProvider,
		gtk.STYLE_PROVIDER_PRIORITY_APPLICATION,
	)

	go func() {
		locations, err := core.GetLocations(authManager.Jwt)

		if err != nil {
			log.Println("An error has occurred")
			return
		}

		sort.Slice(locations.Locations, func(i, j int) bool {
			return locations.Locations[i].Name < locations.Locations[j].Name
		})

		loadingBox.SetVisible(false)
		locationListWrapper.SetVisible(true)

		for i, location := range locations.Locations {
			label := gtk.NewLabel(location.Name)
			locationList.Insert(label, i)
		}

		locationList.Connect("row-selected", func() {
			connectButton.SetSensitive(true)
		})
	}()

	return mainView
}
