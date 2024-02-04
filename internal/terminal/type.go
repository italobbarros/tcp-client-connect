package terminal

import (
	"sync"

	"github.com/gdamore/tcell/v2"
	"github.com/italobbarros/tcp-client-connect/internal/tcp"
	"github.com/rivo/tview"
)

type Terminal struct {
	ManagerConnections *tcp.ManagerConnections
	StatusCh           chan tcp.StatusMsg
	StatusInfoCh       chan tcp.StatusMsg
	Input              chan tcp.DataType
	Output             chan tcp.DataType
	// Terminal
	connection        *tview.TextView
	connectionInfo    *tview.TextView
	sentCommands      *tview.TextView
	receivedResponses *tview.TextView
	app               *tview.Application
	data              *tview.Form
	config            *tview.Form
	stopCh            chan struct{}
	mutex             sync.Mutex
	pages             *tview.Pages
}

type TeminalColors int

const (
	Red TeminalColors = iota
	Yellow
	Green
	Blue
)

// Colors é um mapa que associa valores CustomColor a tcell.Color
var Colors = map[TeminalColors]tcell.Color{
	Red:    tcell.ColorRed,
	Yellow: tcell.ColorYellow,
	Green:  tcell.ColorGreen,
	Blue:   tcell.ColorBlue,
}
