package terminal

import (
	"fmt"
	"time"

	"github.com/italobbarros/tcp-client-connect/internal/tcp"
	"github.com/rivo/tview"
)

func (i *Terminal) Print(msg tcp.DataType, view *tview.TextView) {
	data := fmt.Sprintf("[%s](Conn%d) - %s\n", time.Now().Format("2006-01-02 15:04:05"), msg.ConnId, string(msg.Data))
	view.Write([]byte(data))
	view.ScrollToEnd()
}

func (i *Terminal) PrintStatusConn(value string, color TeminalColors) {
	if i.connection != nil {
		i.app.QueueUpdate(func() {
			//i.connection.Clear()
			i.connection.Write([]byte(value))
			i.connection.SetTextColor(Colors[color])
			i.connection.ScrollToEnd()

		})
	}
}

func (i *Terminal) ClearStatusConn() {
	if i.connection != nil {
		i.app.QueueUpdate(func() {
			i.connection.Clear()
		})
	}
}

func (i *Terminal) PrintStatusInfo(value string, color TeminalColors) {
	if i.connectionInfo != nil {
		i.app.QueueUpdate(func() {
			i.connectionInfo.Clear()
			i.connectionInfo.Write([]byte(value))
			i.connectionInfo.SetTextColor(Colors[color])
		})
	}
}

func (i *Terminal) PrintInput(msg tcp.DataType) {
	if i.sentCommands != nil {
		i.app.QueueUpdate(func() {
			i.Print(msg, i.sentCommands)
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

func (i *Terminal) ClearAll() {
	i.ClearInput()
	//i.ClearOutput()
}
