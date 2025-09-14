package widgets

import (
	_ "embed"
	"math"
	"time"
	"github.com/virtual-designer/urnetwork-client-gnu/core"
	"github.com/diamondburned/gotk4/pkg/gtk/v4"
	"github.com/diamondburned/gotk4/pkg/gdk/v4"
	"github.com/diamondburned/gotk4/pkg/cairo"
)

//go:embed connectedView.ui
var connectedViewUiXML string

//go:embed connectedView.css
var connectedViewCSS string

type ConnectedView struct {
	*gtk.Box
	stack *gtk.Stack

	connectInfoWrapper *gtk.Box
	connectCircleWrapper *gtk.Box
	connectStatus *gtk.Label
	connectLocation *gtk.Label
	disconnectButton *gtk.Button
	connectCircle *gtk.DrawingArea
}

func (connectedView *ConnectedView) DrawCircle(r float64, g float64, b float64) {
	drawingArea := connectedView.connectCircle

	drawingArea.SetDrawFunc(func (area *gtk.DrawingArea, context *cairo.Context, width, height int) {
		area.SetContentWidth(140)
		area.SetContentHeight(140)
		area.SetHAlign(gtk.AlignCenter)

		context.SetSourceRGBA(r / 255.0, g / 255.0, b / 255.0, 0.4)
		context.Arc(70, 70, 70, 0, 2 * math.Pi)
		context.Fill()
		context.SetSourceRGBA(r / 255.0, g / 255.0, b / 255.0, 0.9)
		context.Arc(70, 70, 50, 0, 2 * math.Pi)
		context.Fill()
		context.SetSourceRGBA(r / 255.0, g / 255.0, b / 255.0, 0.5)
		context.Arc(70, 70, 60, 0, 2 * math.Pi)
		context.Fill()
	})
}

func (connectedView *ConnectedView) OnConnect(location *core.APILocationResult) {
	connectedView.DrawCircle(240.0, 224.0, 0.0)
	connectedView.connectStatus.SetLabel("Connecting")
	connectedView.connectLocation.SetLabel(location.Name)
	time.Sleep(2 * time.Second)
	connectedView.DrawCircle(52.0, 235.0, 85.0)
	connectedView.connectStatus.SetCSSClasses([]string {"connect_status_connected"})
	connectedView.connectStatus.SetLabel("Connected")
	connectedView.disconnectButton.SetVisible(true)
	connectedView.disconnectButton.SetSensitive(true)
}

func (connectedView *ConnectedView) OnDisconnect() {
	connectedView.DrawCircle(240.0, 224.0, 0.0)
	connectedView.connectStatus.SetCSSClasses([]string {"connect_status"})
	connectedView.connectStatus.SetLabel("Disconnecting")
	connectedView.disconnectButton.SetSensitive(false)
	connectedView.stack.SetVisibleChildName("MainView")
}

func NewConnectedView(stack *gtk.Stack) *ConnectedView {
	builder := gtk.NewBuilderFromString(connectedViewUiXML)

	cssProvider := gtk.NewCSSProvider()
	cssProvider.LoadFromData(connectedViewCSS)

	connectedView := & ConnectedView {
		Box: builder.GetObject("ConnectedView").Cast().(*gtk.Box),
		stack: stack,
		connectInfoWrapper: builder.GetObject("ConnectInfoWrapper").Cast().(*gtk.Box),
		connectCircleWrapper: builder.GetObject("ConnectCircleWrapper").Cast().(*gtk.Box),
		connectCircle: builder.GetObject("ConnectCircle").Cast().(*gtk.DrawingArea),
		connectStatus: builder.GetObject("ConnectStatusLabel").Cast().(*gtk.Label),
		connectLocation: builder.GetObject("ConnectLocationLabel").Cast().(*gtk.Label),
		disconnectButton: builder.GetObject("DisconnectButton").Cast().(*gtk.Button),
	}

	gtk.StyleContextAddProviderForDisplay(
		gdk.DisplayGetDefault(),
		cssProvider,
		gtk.STYLE_PROVIDER_PRIORITY_APPLICATION,
	)

	connectedView.disconnectButton.Connect("clicked", func () {
		connectedView.OnDisconnect()
	})

	return connectedView
}
