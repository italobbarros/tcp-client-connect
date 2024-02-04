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
			data := Colors[color]
			data = append(data, []byte(value)...)
			i.connection.Write(data)
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
			data := Colors[color]
			data = append(data, []byte(value)...)
			i.connectionInfo.Write(data)
		})
	}
}

func (i *Terminal) PrintInput(msg tcp.DataType) {
	if i.outputView != nil {
		i.app.QueueUpdate(func() {
			i.Print(msg, i.outputView)
		})
	}
}

func (i *Terminal) ClearInput() {
	if i.outputView != nil {
		i.app.QueueUpdate(func() {
			i.outputView.Clear()
		})
	}
}

func (i *Terminal) ClearOutput() {
	if i.inputView != nil {
		i.app.QueueUpdate(func() {
			i.inputView.Clear()
		})
	}
}

func (i *Terminal) ClearAll() {
	i.ClearInput()
	//i.ClearOutput()
}
