package terminal

import (
	"sync"

	"github.com/gdamore/tcell/v2"
	"github.com/italobbarros/tcp-client-connect/internal/client"
	"github.com/rivo/tview"
)

type Terminal struct {
	managerClient *client.ManagerClients
	StatusCh      chan client.StatusMsg
	Input         chan client.DataType
	Output        chan client.DataType
	// Terminal
	connection        *tview.TextView
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
