package widgets

import (
	_ "embed"
	"math"
	"os"
	"os/signal"
	"os/exec"
	"syscall"
	"fmt"
	"errors"
	"context"
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
	authManager *core.AuthManager
	stack *gtk.Stack

	connectInfoWrapper *gtk.Box
	connectCircleWrapper *gtk.Box
	connectStatus *gtk.Label
	connectLocation *gtk.Label
	disconnectButton *gtk.Button
	retryButton *gtk.Button
	backButton *gtk.Button
	connectCircleDrawingArea *gtk.DrawingArea

	lastLocation *core.APILocationResult
	clientProcess *exec.Cmd
	currentCancel context.CancelFunc
	isVisible bool
}

func (connectedView *ConnectedView) DrawCircle(r float64, g float64, b float64) {
	drawFunc := func (area *gtk.DrawingArea, context *cairo.Context, width, height int) {
		context.SetSourceRGBA(r / 255.0, g / 255.0, b / 255.0, 0.4)
		context.Arc(70, 70, 70, 0, 2 * math.Pi)
		context.Fill()
		context.SetSourceRGBA(r / 255.0, g / 255.0, b / 255.0, 0.9)
		context.Arc(70, 70, 50, 0, 2 * math.Pi)
		context.Fill()
		context.SetSourceRGBA(r / 255.0, g / 255.0, b / 255.0, 0.5)
		context.Arc(70, 70, 60, 0, 2 * math.Pi)
		context.Fill()
	}

	connectedView.connectCircleDrawingArea.SetDrawFunc(drawFunc)
	connectedView.connectCircleDrawingArea.QueueDraw()
}

func (connectedView *ConnectedView) OnFail(err error) {
	connectedView.DrawCircle(252.0, 78.0, 3.0)
	connectedView.connectStatus.SetLabel("Unable to connect")
	connectedView.connectLocation.SetLabel("")
	connectedView.retryButton.SetVisible(true)
	connectedView.backButton.SetVisible(true)
	connectedView.connectStatus.SetCSSClasses([]string {"connect_status_error"})
}

func (connectedView *ConnectedView) OnConnect(location *core.APILocationResult) {
	connectedView.isVisible = true

	if connectedView.currentCancel != nil {
		connectedView.currentCancel()
	}

	ctx, cancel := context.WithCancel(context.Background())
	connectedView.currentCancel = cancel
	connectedView.lastLocation = location

	connectedView.DrawCircle(240.0, 224.0, 0.0)
	connectedView.connectStatus.SetLabel("Connecting")
	connectedView.connectLocation.SetLabel(location.Name)
	connectedView.retryButton.SetVisible(false)
	connectedView.backButton.SetVisible(false)
	connectedView.connectStatus.SetCSSClasses([]string {"connect_status"})

	urnetworkClientPath := os.Args[0] + "/bin/urnetwork-client"

	if envClientPath := os.Getenv("URNETWORK_CLIENT_PATH"); envClientPath != "" {
		urnetworkClientPath = envClientPath
	}

	cmd := exec.Command("pkexec", "--keep-cwd", urnetworkClientPath, "quick-connect",
		"--jwt", connectedView.authManager.Jwt,
		"--tun=tun0", "--stats_interval", "0", "--default_route",
		"--location_query=country:" + location.Name, "--dns_bootstrap=cache")

	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Start(); err != nil {
		connectedView.OnFail(err)
		return
	}

	sigChannel := make(chan os.Signal, 1)
	signal.Notify(sigChannel, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		sig := <-sigChannel
		connectedView.clientProcess = nil

		if sig == syscall.SIGKILL {
			return
		}

		fmt.Println("Signal received: ", sig)
		exec.Command("pkexec", "kill", "-TERM", fmt.Sprintf("%d", cmd.Process.Pid)).Run()
		os.Exit(0)
	}();

	go func() {
		done := make(chan error, 1)

		go func() {
			done <- cmd.Wait()
		}()

		select {
		case <-ctx.Done():
			fmt.Println("Terminated")

			if cmd.Process != nil {
				exec.Command("pkexec", "kill", "-TERM", fmt.Sprintf("%d", cmd.Process.Pid)).Run()
				connectedView.clientProcess = nil
			}

		case err := <-done:
			if err == nil {
				err = errors.New("Unknown error occurred")
			}

			sigChannel <- syscall.SIGKILL
			connectedView.clientProcess = nil

			if connectedView.isVisible {
				connectedView.OnFail(err)
				connectedView.disconnectButton.SetVisible(false)
				connectedView.disconnectButton.SetSensitive(false)
			}
		}
	}();

	connectedView.clientProcess	= cmd

	connectedView.DrawCircle(52.0, 235.0, 85.0)
	connectedView.connectStatus.SetCSSClasses([]string {"connect_status_connected"})
	connectedView.connectStatus.SetLabel("Connected")
	connectedView.disconnectButton.SetVisible(true)
	connectedView.disconnectButton.SetSensitive(true)
}

func (connectedView *ConnectedView) OnDisconnect() {
	if connectedView.clientProcess != nil {
		err := exec.Command("pkexec", "kill", "-TERM", fmt.Sprintf("%d", connectedView.clientProcess.Process.Pid)).Run()

		if err != nil {
			return
		}

		connectedView.clientProcess = nil
	}

	connectedView.DrawCircle(240.0, 224.0, 0.0)
	connectedView.connectStatus.SetCSSClasses([]string {"connect_status"})
	connectedView.connectStatus.SetLabel("Disconnecting")
	connectedView.disconnectButton.SetSensitive(false)
	connectedView.backButton.SetVisible(false)
	connectedView.stack.SetVisibleChildName("MainView")
}

func NewConnectedView(authManager *core.AuthManager, stack *gtk.Stack) *ConnectedView {
	builder := gtk.NewBuilderFromString(connectedViewUiXML)

	cssProvider := gtk.NewCSSProvider()
	cssProvider.LoadFromData(connectedViewCSS)

	connectedView := & ConnectedView {
		Box: builder.GetObject("ConnectedView").Cast().(*gtk.Box),
		authManager: authManager,
		stack: stack,
		connectInfoWrapper: builder.GetObject("ConnectInfoWrapper").Cast().(*gtk.Box),
		connectCircleWrapper: builder.GetObject("ConnectCircleWrapper").Cast().(*gtk.Box),
		connectCircleDrawingArea: builder.GetObject("ConnectCircle").Cast().(*gtk.DrawingArea),
		connectStatus: builder.GetObject("ConnectStatusLabel").Cast().(*gtk.Label),
		connectLocation: builder.GetObject("ConnectLocationLabel").Cast().(*gtk.Label),
		disconnectButton: builder.GetObject("DisconnectButton").Cast().(*gtk.Button),
		retryButton: builder.GetObject("RetryButton").Cast().(*gtk.Button),
		backButton: builder.GetObject("BackButton").Cast().(*gtk.Button),
		isVisible: false,
	}

	gtk.StyleContextAddProviderForDisplay(
		gdk.DisplayGetDefault(),
		cssProvider,
		gtk.STYLE_PROVIDER_PRIORITY_APPLICATION,
	)

	connectedView.disconnectButton.Connect("clicked", func () {
		fmt.Println("Disconnecting")
		connectedView.OnDisconnect()
		connectedView.isVisible = false
	})

	connectedView.backButton.Connect("clicked", func () {
		connectedView.OnDisconnect()
		connectedView.isVisible = false
	})

	connectedView.retryButton.Connect("clicked", func () {
		if connectedView.lastLocation != nil {
			connectedView.OnConnect(connectedView.lastLocation)
		}
	})

	return connectedView
}
