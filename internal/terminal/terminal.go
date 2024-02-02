package terminal

import (
	"strconv"
	"strings"
	"time"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

func NewTerminal(serverCommandCh *chan string, userCommandCh *chan string) *Terminal {
	return &Terminal{
		serverCommandCh: serverCommandCh,
		userCommandCh:   userCommandCh,
		stopCh:          make(chan struct{}),
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
		SetDynamicColors(true).
		SetWrap(true)
	t.connection.SetBorder(true).SetTitle("Status").SetTitleAlign(tview.AlignLeft)

	// Cria dois novos TextViews para os comandos enviados e recebidos
	t.sentCommands = tview.NewTextView().
		SetDynamicColors(true).
		SetWrap(true)
	t.sentCommands.SetBorder(true).SetTitle("Enviados").SetTitleAlign(tview.AlignLeft)
	t.receivedResponses = tview.NewTextView().
		SetDynamicColors(true).
		SetWrap(true)
	t.receivedResponses.SetBorder(true).SetTitle("Recebidos").SetTitleAlign(tview.AlignLeft)

	// Cria um novo Form para a entrada do usuário
	t.data = tview.NewForm().
		AddInputField("Dados", "", 50, nil, nil) //.
	t.data.SetBorder(true).SetTitle("Enviar dados").SetTitleAlign(tview.AlignLeft)

	t.config = tview.NewForm().
		AddDropDown("Intervalo", []string{"None", "5s", "10s", "60s"}, 0, nil).
		AddButton("Save", nil)
	t.config.SetBorder(true).SetTitle("Config param").SetTitleAlign(tview.AlignLeft)

	t.config.GetButton(0).SetSelectedFunc(func() {
		_, interval := t.config.GetFormItemByLabel("Intervalo").(*tview.DropDown).GetCurrentOption()
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
						command := t.data.GetFormItemByLabel("Dados").(*tview.InputField).GetText()
						if len(command) == 0 {
							continue
						}
						*t.userCommandCh <- command
						t.Print(command, t.sentCommands)
						t.app.SetFocus(t.data.GetFormItemByLabel("Dados"))
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
	})

	t.data.GetFormItemByLabel("Dados").(*tview.InputField).
		SetDoneFunc(func(key tcell.Key) {
			if key == tcell.KeyEnter {
				// Chama a função de envio ao pressionar "Enter"
				command := t.data.GetFormItemByLabel("Dados").(*tview.InputField).GetText()
				if len(command) > 0 {
					*t.userCommandCh <- command
					t.Print(command, t.sentCommands)
					t.data.GetFormItemByLabel("Dados").(*tview.InputField).SetText("")
					t.app.SetFocus(t.data.GetFormItemByLabel("Dados"))
				}
			}
		})

	background := tview.NewFlex().
		AddItem(tview.NewFlex().SetDirection(tview.FlexRow).
			AddItem(t.connection, 4, 1, false).
			AddItem(tview.NewFlex().SetDirection(tview.FlexColumn).
				AddItem(t.sentCommands, 0, 1, false).
				AddItem(t.receivedResponses, 0, 1, false), 0, 2, false).
			AddItem(tview.NewFlex().SetDirection(tview.FlexColumn).
				AddItem(t.data, 0, 1, true).
				AddItem(t.config, 0, 1, false), 7, 1, false,
			), 0, 1, true)

	modal := t.ConfigModal(endCh)

	t.pages = tview.NewPages().
		AddPage("background", background, true, true).
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
		case response := <-*t.serverCommandCh:
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
