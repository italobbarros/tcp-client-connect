package terminal

import (
	"strconv"
	"strings"
	"time"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

func NewTerminal(serverCommandCh chan []byte, userCommandCh chan []byte, clientIsConnected func() bool) *Terminal {
	return &Terminal{
		serverCommandCh:   serverCommandCh,
		userCommandCh:     userCommandCh,
		stopCh:            make(chan struct{}),
		clientIsConnected: clientIsConnected,
	}
}

func (t *Terminal) Create(endCh chan struct{}) {
	t.app = tview.NewApplication()
	go func() {
		for {
			select {
			case <-endCh:
				t.app.Stop()
			case <-time.After(time.Duration(1) * time.Second):
				t.app.Draw()
			}
		}
	}()
	go func() {
		for {
			select {
			case <-endCh:
				return
			case <-time.After(time.Duration(5) * time.Minute):
				t.ClearInput()
				t.ClearOutput()
			}
		}
	}()
	t.connection = tview.NewTextView().
		SetText("Desconectado!").
		SetTextColor(tcell.ColorRed).
		SetDynamicColors(true).
		SetWrap(true).
		SetTextAlign(tview.AlignCenter)
	t.connection.SetBorder(true).SetTitle("Status").SetTitleAlign(tview.AlignLeft)

	// Cria dois novos TextViews para os comandos enviados e recebidos
	t.sentCommands = tview.NewTextView().
		SetDynamicColors(true).
		SetWrap(true)
	t.sentCommands.SetBorder(true).SetTitle("Output").SetTitleAlign(tview.AlignLeft)
	t.receivedResponses = tview.NewTextView().
		SetDynamicColors(true).
		SetWrap(true)
	t.receivedResponses.SetBorder(true).SetTitle("Input").SetTitleAlign(tview.AlignLeft)

	// Cria um novo Form para a entrada do usuário
	t.data = tview.NewForm().
		AddInputField("Data", "", 150, nil, nil) //.
	t.data.SetBorder(true).SetTitle("Send Data").SetTitleAlign(tview.AlignLeft)

	t.config = tview.NewForm().
		AddDropDown("Time", []string{"None", "5s", "10s", "60s"}, 0, nil).
		AddDropDown("View", []string{"Input", "All"}, 0, nil).
		AddButton("Save", func() {
			_, interval := t.config.GetFormItemByLabel("Time").(*tview.DropDown).GetCurrentOption()
			if interval != "None" {
				close(t.stopCh)
				go func() {
					time.Sleep(1 * time.Millisecond)
					t.stopCh = make(chan struct{})
					var interValue int = 0
					for {
						select {
						case <-endCh:
							return
						case <-t.stopCh:
							return
						case <-time.After(time.Duration(interValue) * time.Second):
							if !t.clientIsConnected() {
								interValue = 1
								continue
							}
							command := t.data.GetFormItemByLabel("Data").(*tview.InputField).GetText()
							if len(command) == 0 {
								continue
							}
							data := []byte(command)
							t.userCommandCh <- data
							t.Print(data, t.sentCommands)
							//t.app.SetFocus(t.data.GetFormItemByLabel("Data"))
							numberStr := strings.TrimSuffix(interval, "s")
							v, err := strconv.Atoi(numberStr)
							if err != nil {
								return
							}
							interValue = v
						}
					}
				}()
			}
			_, pages := t.config.GetFormItemByLabel("View").(*tview.DropDown).GetCurrentOption()
			if pages == "All" {
				if t.pages != nil {
					t.pages.ShowPage("page_in_and_out")
					t.pages.HidePage("page_out_only")
				}
			}
			if pages == "Input" {
				if t.pages != nil {
					t.pages.ShowPage("page_out_only")
					t.pages.HidePage("page_in_and_out")
				}
			}
		})
	t.config.SetBorder(true).SetTitle("Config").SetTitleAlign(tview.AlignLeft)
	t.data.GetFormItemByLabel("Data").(*tview.InputField).
		SetDoneFunc(func(key tcell.Key) {
			if !t.clientIsConnected() {
				return
			}
			if key == tcell.KeyEnter {
				// Chama a função de envio ao pressionar "Enter"
				command := t.data.GetFormItemByLabel("Data").(*tview.InputField).GetText()
				if len(command) == 0 {
					return
				}
				data := []byte(command)
				t.userCommandCh <- data
				t.Print(data, t.sentCommands)
				t.data.GetFormItemByLabel("Data").(*tview.InputField).SetText("")
				t.app.SetFocus(t.data.GetFormItemByLabel("Data"))
			}
		})

	page_in_and_out := tview.NewFlex().
		AddItem(tview.NewFlex().SetDirection(tview.FlexRow).
			AddItem(t.connection, 3, 1, false).
			AddItem(tview.NewFlex().SetDirection(tview.FlexColumn).
				AddItem(t.sentCommands, 0, 1, false).
				AddItem(t.receivedResponses, 0, 1, false).
				AddItem(t.config, 16, 1, false), 0, 1, false).
			AddItem(t.data, 5, 1, true), 0, 1, true)

	page_out_only := tview.NewFlex().
		AddItem(tview.NewFlex().SetDirection(tview.FlexRow).
			AddItem(t.connection, 3, 1, false).
			AddItem(tview.NewFlex().SetDirection(tview.FlexColumn).
				AddItem(t.receivedResponses, 0, 1, false).
				AddItem(t.config, 16, 1, false), 0, 1, false).
			AddItem(t.data, 5, 1, true), 0, 1, true)

	modal := t.ConfigModal(endCh)

	t.pages = tview.NewPages().
		AddPage("page_in_and_out", page_in_and_out, true, false).
		AddPage("page_out_only", page_out_only, true, true).
		AddPage("modal", modal, false, false)
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
		case response := <-t.serverCommandCh:
			// Imprime a resposta recebida
			if t.receivedResponses != nil {
				t.Print(response, t.receivedResponses)
				t.app.Draw()
			}
		}
	}
}

func (t *Terminal) ConfigModal(endCh chan struct{}) tview.Primitive {
	flexModal := func(p tview.Primitive, width, height int) tview.Primitive {
		return tview.NewFlex().
			AddItem(nil, 0, 1, false).
			AddItem(tview.NewFlex().SetDirection(tview.FlexRow).
				AddItem(nil, 0, 1, false).
				AddItem(p, height, 1, true).
				AddItem(nil, 0, 1, false), width, 1, true).
			AddItem(nil, 0, 1, false)
	}
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

	return flexModal(modalBase, 10, 10)
}
