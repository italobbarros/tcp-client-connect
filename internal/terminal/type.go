package terminal

import (
	"sync"

	"github.com/italobbarros/tcp-client-connect/internal/tcp"
	"github.com/rivo/tview"
)

type Terminal struct {
	ManagerConnections *tcp.ManagerConnections
	StatusCh           chan tcp.StatusMsg
	StatusInfoCh       chan tcp.StatusMsg
	Input              chan tcp.DataType
	//Output             chan tcp.DataType
	// Terminal
	connection     *tview.TextView
	connectionInfo *tview.TextView
	outputView     *tview.TextView
	inputView      *tview.TextView
	app            *tview.Application
	data           *tview.Form
	config         *tview.Form
	timerCh        chan struct{}
	closingTimer   sync.Once
	loopback       bool
	mutex          sync.Mutex
	pages          *tview.Pages
}

type TeminalColors int

const (
	Red TeminalColors = iota
	Yellow
	Green
	Blue
)

// Colors é um mapa que associa valores CustomColor a tcell.Color
var Colors = map[TeminalColors][]byte{
	Red:    []byte("[red]"),
	Yellow: []byte("[yellow]"),
	Green:  []byte("[green]"),
	Blue:   []byte("[blue]"),
}
