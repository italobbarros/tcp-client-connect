package terminal

import (
	"sync"

	"github.com/rivo/tview"
)

type Terminal struct {
	serverCommandCh   *chan string
	userCommandCh     *chan string
	sentCommands      *tview.TextView
	receivedResponses *tview.TextView
	app               *tview.Application
	data              *tview.Form
	config            *tview.Form
	stopCh            chan struct{}
	mutex             sync.Mutex
}
