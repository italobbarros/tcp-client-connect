package terminal

import (
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/gdamore/tcell/v2"
	"github.com/italobbarros/tcp-client-connect/internal/tcp"
	"github.com/rivo/tview"
)

func NewTerminal(managerConnections *tcp.ManagerConnections) *Terminal {
	totalConnections := managerConnections.GetNumberConnections()
	return &Terminal{
		ManagerConnections: managerConnections,
		timerCh:            make(chan struct{}),
		StatusInfoCh:       make(chan tcp.StatusMsg),
		StatusCh:           make(chan tcp.StatusMsg, totalConnections),
		Input:              make(chan tcp.DataType, totalConnections),
		loopback:           false,
	}
}

func (t *Terminal) Create(endCh chan struct{}) {
	t.app = tview.NewApplication()
	go func() {
		for {
			select {
			case <-endCh:
				t.app.Stop()
			case stats := <-t.StatusCh:
				t.PrintStatusConn(stats.Msg, TeminalColors(stats.Color))
			case stats := <-t.StatusInfoCh:
				t.PrintStatusInfo(stats.Msg, TeminalColors(stats.Color))
			case <-time.After(time.Duration(1) * time.Second):
				t.app.Draw()
			}
		}
	}()
	t.connection = tview.NewTextView().
		SetText("Desconectado!").
		SetTextColor(tcell.ColorRed).
		SetDynamicColors(true).
		SetWrap(true).
		SetTextAlign(tview.AlignCenter).Clear().
		SetDoneFunc(func(key tcell.Key) {
			if key == tcell.KeyEnter {
				t.pages.ShowPage("StatusConnection")
			}
		}).SetToggleHighlights(true).SetMaxLines(100)

	t.connection.SetBorder(true).SetTitle("Status").SetTitleAlign(tview.AlignLeft)

	t.connectionInfo = tview.NewTextView().
		SetDynamicColors(true).
		SetWrap(true).
		SetTextAlign(tview.AlignCenter).SetMaxLines(1)
	t.connectionInfo.SetBorder(true).SetTitle("Info").SetTitleAlign(tview.AlignLeft)

	// Cria dois novos TextViews para os comandos enviados e recebidos
	t.outputView = tview.NewTextView().
		SetDynamicColors(true).
		SetWrap(true).SetMaxLines(100)
	t.outputView.SetBorder(true).SetTitle("Output").SetTitleAlign(tview.AlignLeft)
	t.inputView = tview.NewTextView().
		SetDynamicColors(true).
		SetWrap(true).SetMaxLines(100)
	t.inputView.SetBorder(true).SetTitle("Input").SetTitleAlign(tview.AlignLeft)

	// Cria um novo Form para a entrada do usuário
	t.data = tview.NewForm().
		AddInputField("Data", "", 150, nil, nil) //.
	t.data.SetBorder(true).SetTitle("Send Data").SetTitleAlign(tview.AlignLeft)

	t.config = tview.NewForm().
		AddDropDown("Time", []string{"None", "5s", "10s", "60s"}, 0, func(option string, optionIndex int) {
			t.CloseTimer()
			if option == "None" {
				return
			}
			go func() {
				time.Sleep(1 * time.Millisecond)
				t.timerCh = make(chan struct{})
				t.closingTimer = sync.Once{}
				var interValue int = 0
				for {
					select {
					case <-endCh:
						return
					case <-t.timerCh:
						return
					case <-time.After(time.Duration(interValue) * time.Second):
						command := t.data.GetFormItemByLabel("Data").(*tview.InputField).GetText()
						if len(command) == 0 {
							continue
						}
						data := tcp.DataType{
							Data:   []byte(command),
							ConnId: 0,
						}
						clientsId := t.ManagerConnections.SendDataToConnections(data)
						go func() {
							for id := range clientsId {
								data.ConnId = id
								t.Print(data, t.outputView)
							}
						}()
						numberStr := strings.TrimSuffix(option, "s")
						v, err := strconv.Atoi(numberStr)
						if err != nil {
							return
						}
						interValue = v
					}
				}
			}()

		}).
		AddDropDown("View", []string{"Input", "All"}, 0, func(option string, optionIndex int) {
			if option == "All" {
				if t.pages != nil {
					t.pages.ShowPage("page_in_and_out")
					t.pages.HidePage("page_in_only")
				}
			}
			if option == "Input" {
				if t.pages != nil {
					t.pages.ShowPage("page_in_only")
					t.pages.HidePage("page_in_and_out")
				}
			}
		}).
		AddCheckbox("Loopback", false, func(checked bool) {
			if checked {
				t.loopback = true
			} else {
				t.loopback = false
			}
		})

	t.config.SetBorder(true).SetTitle("Config").SetTitleAlign(tview.AlignLeft)
	t.data.GetFormItemByLabel("Data").(*tview.InputField).
		SetDoneFunc(func(key tcell.Key) {
			if key == tcell.KeyEnter {
				// Chama a função de envio ao pressionar "Enter"
				command := t.data.GetFormItemByLabel("Data").(*tview.InputField).GetText()
				if len(command) == 0 {
					return
				}
				data := tcp.DataType{
					Data:   []byte(command),
					ConnId: 0,
				}
				go func() {
					clientsId := t.ManagerConnections.SendDataToConnections(data)
					for id := range clientsId {
						data.ConnId = id
						t.Print(data, t.outputView)
					}
				}()
				t.data.GetFormItemByLabel("Data").(*tview.InputField).SetText("")
				t.app.SetFocus(t.data.GetFormItemByLabel("Data"))
			}
		})

	page_in_and_out := tview.NewFlex().
		AddItem(tview.NewFlex().SetDirection(tview.FlexRow).
			AddItem(tview.NewFlex().SetDirection(tview.FlexColumn).
				AddItem(t.connection, 0, 1, false).
				AddItem(t.connectionInfo, 22, 1, false), 3, 1, false).
			AddItem(tview.NewFlex().SetDirection(tview.FlexColumn).
				AddItem(tview.NewFlex().SetDirection(tview.FlexRow).
					AddItem(tview.NewFlex().SetDirection(tview.FlexColumn).
						AddItem(t.outputView, 0, 1, false).
						AddItem(t.inputView, 0, 1, false), 0, 1, true).
					AddItem(t.data, 5, 1, true), 0, 1, true).
				AddItem(t.config, 22, 1, false), 0, 1, false), 0, 1, true)

	page_in_only := tview.NewFlex().
		AddItem(tview.NewFlex().SetDirection(tview.FlexRow).
			AddItem(tview.NewFlex().SetDirection(tview.FlexColumn).
				AddItem(t.connection, 0, 1, false).
				AddItem(t.connectionInfo, 22, 1, false), 3, 1, false).
			AddItem(tview.NewFlex().SetDirection(tview.FlexColumn).
				AddItem(tview.NewFlex().SetDirection(tview.FlexRow).
					AddItem(t.inputView, 0, 1, false).
					AddItem(t.data, 5, 1, true), 0, 1, true).
				AddItem(t.config, 22, 1, false), 0, 1, false), 0, 1, true)
		//AddItem(t.data, 5, 1, true), 0, 1, true)

	modalPage := t.ConfigModal(endCh)
	StatusConnectionPage := t.ExpandingStatusModal()

	t.pages = tview.NewPages().
		AddPage("page_in_and_out", page_in_and_out, true, false).
		AddPage("page_in_only", page_in_only, true, true).
		AddPage("modal", modalPage, false, false).
		AddPage("StatusConnection", StatusConnectionPage, true, false)
	// Define o flex como a raiz da aplicação
	if err := t.app.SetRoot(t.pages, true).EnableMouse(true).Run(); err != nil {
		panic(err)
	}
}

func (t *Terminal) ListenServerResponse(endCh chan struct{}) {
	for {
		select {
		case <-endCh:
			return
		case response := <-t.Input:
			// Imprime a resposta recebida
			if t.inputView != nil {
				t.Print(response, t.inputView)
				t.app.Draw()
			}
			if t.loopback {
				if len(response.Data) == 0 {
					continue
				}
				t.ManagerConnections.SendDataToConnections(response)
				go t.Print(response, t.outputView)
				time.Sleep(time.Microsecond * 1)
			}
		}
	}
}

func (t *Terminal) CloseTimer() {
	t.closingTimer.Do(func() {
		close(t.timerCh)
	})
}

func (t *Terminal) FlexModalPrimitive(p tview.Primitive, width, height int) tview.Primitive {
	return tview.NewFlex().
		AddItem(nil, 0, 1, false).
		AddItem(tview.NewFlex().SetDirection(tview.FlexRow).
			AddItem(nil, 0, 1, false).
			AddItem(p, height, 1, true).
			AddItem(nil, 0, 1, false), width, 1, true).
		AddItem(nil, 0, 1, false)

}

func (t *Terminal) ConfigModal(endCh chan struct{}) tview.Primitive {

	modalBase := tview.NewModal().
		SetText("Do you want to quit the application?").
		AddButtons([]string{"Quit", "Cancel"}).
		SetButtonActivatedStyle(tcell.StyleDefault.
			Foreground(tcell.ColorRed).
			Background(tcell.ColorWhite).
			Bold(true)).
		SetDoneFunc(func(buttonIndex int, buttonLabel string) {
			if buttonLabel == "Quit" {
				t.app.Stop()
				close(endCh)
			} else {
				if t.pages != nil {
					t.pages.HidePage("modal")
				}
			}
		})
	t.app.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyCtrlC {
			if t.pages != nil {
				t.pages.ShowPage("modal")
			}
			return nil
		}
		return event
	})

	return t.FlexModalPrimitive(modalBase, 10, 10)
}

func (t *Terminal) ExpandingStatusModal() tview.Primitive {
	statusFlex := tview.NewFlex().
		AddItem(tview.NewFlex().SetDirection(tview.FlexRow).
			AddItem(t.connection.SetTextAlign(tview.AlignLeft), 0, 1, false).
			AddItem(tview.NewButton("Exit").SetLabelColor(tcell.ColorRed).
				SetSelectedFunc(func() {
					t.pages.HidePage("StatusConnection")
				}), 3, 1, true), 0, 1, true)
	return statusFlex
}
