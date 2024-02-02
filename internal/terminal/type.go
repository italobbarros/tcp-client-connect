package terminal

import (
	"sync"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type Terminal struct {
	serverCommandCh   *chan string
	userCommandCh     *chan string
	connection        *tview.TextView
	sentCommands      *tview.TextView
	receivedResponses *tview.TextView
	app               *tview.Application
	data              *tview.Form
	config            *tview.Form
	stopCh            chan struct{}
	mutex             sync.Mutex
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
