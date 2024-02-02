package terminal

import (
	"time"

	"github.com/rivo/tview"
)

func (i *Terminal) Print(value string, view *tview.TextView) {
	data := time.Now().Format("2006-01-02 15:04:05") + " - " + value + "\n"
	view.Write([]byte(data))
	view.ScrollToEnd()
}

func (i *Terminal) PrintStatus(value string, color TeminalColors) {
	if i.connection != nil {
		i.app.QueueUpdate(func() {
			i.connection.Clear()
			i.connection.Write([]byte(value))
			i.connection.SetTextColor(Colors[color])
			i.connection.ScrollToEnd()

		})
	}
}

func (i *Terminal) PrintInput(value string) {
	if i.sentCommands != nil {
		i.app.QueueUpdate(func() {
			i.Print(value, i.sentCommands)
		})
	}
}

func (i *Terminal) ClearInput() {
	if i.sentCommands != nil {
		i.app.QueueUpdate(func() {
			i.sentCommands.Clear()
		})
	}
}

func (i *Terminal) ClearOutput() {
	if i.receivedResponses != nil {
		i.app.QueueUpdate(func() {
			i.receivedResponses.Clear()
		})
	}
}
