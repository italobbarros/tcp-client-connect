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

func (i *Terminal) Create(doneCh chan struct{}) {
	i.app = tview.NewApplication()
	go func() {
		for {
			select {
			case <-doneCh:
				i.app.Stop()
			case <-time.After(time.Duration(1) * time.Second):
				i.app.Draw()
			}
		}
	}()
	go func() {
		for {
			// Use tview.App.QueueUpdate para garantir operações seguras na GUI
			if i.sentCommands != nil {
				i.app.QueueUpdate(func() {
					i.sentCommands.Clear()
				})
			}
			if i.receivedResponses != nil {
				i.app.QueueUpdate(func() {
					i.receivedResponses.Clear()
				})
			}
			time.Sleep(5 * time.Minute)
		}
	}()
	// Cria dois novos TextViews para os comandos enviados e recebidos
	i.sentCommands = tview.NewTextView().
		SetDynamicColors(true).
		SetWrap(true)
	i.sentCommands.SetBorder(true).SetTitle("Enviados").SetTitleAlign(tview.AlignLeft)
	i.receivedResponses = tview.NewTextView().
		SetDynamicColors(true).
		SetWrap(true)
	i.receivedResponses.SetBorder(true).SetTitle("Recebidos").SetTitleAlign(tview.AlignLeft)

	// Cria um novo Form para a entrada do usuário
	i.data = tview.NewForm().
		AddInputField("Dados", "", 100, nil, nil) //.
	i.data.SetBorder(true).SetTitle("Enviar dados").SetTitleAlign(tview.AlignLeft)

	i.config = tview.NewForm().
		AddDropDown("Intervalo", []string{"None", "5s", "10s", "60s"}, 0, nil).
		AddButton("Save", nil)
	i.config.SetBorder(true).SetTitle("Config param").SetTitleAlign(tview.AlignLeft)

	i.config.GetButton(0).SetSelectedFunc(func() {
		_, interval := i.config.GetFormItemByLabel("Intervalo").(*tview.DropDown).GetCurrentOption()
		if interval != "None" {
			close(i.stopCh)
			go func() {
				time.Sleep(1 * time.Millisecond)
				i.stopCh = make(chan struct{})
				var interValue int = 0
				for {
					select {
					case <-i.stopCh:
						return
					case <-time.After(time.Duration(interValue) * time.Second):
						i.mutex.Lock()
						command := i.data.GetFormItemByLabel("Dados").(*tview.InputField).GetText()
						*i.userCommandCh <- command
						i.Print(command, i.sentCommands)
						i.app.SetFocus(i.config.GetFormItemByLabel("Dados"))
						i.mutex.Unlock()
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

	i.data.GetFormItemByLabel("Dados").(*tview.InputField).
		SetDoneFunc(func(key tcell.Key) {
			if key == tcell.KeyEnter {
				// Chama a função de envio ao pressionar "Enter"
				i.mutex.Lock()
				command := i.data.GetFormItemByLabel("Dados").(*tview.InputField).GetText()
				if len(command) > 0 {
					*i.userCommandCh <- command
					i.Print(command, i.sentCommands)
					i.data.GetFormItemByLabel("Dados").(*tview.InputField).SetText("")
					i.app.SetFocus(i.data.GetFormItemByLabel("Dados"))
				}
				i.mutex.Unlock()
			}
		})
	// Cria um novo Grid e adiciona o Form, os TextViews e um espaço vazio
	grid := tview.NewGrid().
		AddItem(i.sentCommands, 0, 0, 3, 1, 0, 0, false).
		AddItem(i.receivedResponses, 0, 1, 3, 1, 0, 0, false).
		AddItem(i.data, 3, 0, 1, 2, 0, 0, true).
		AddItem(i.config, 4, 0, 1, 2, 0, 0, true)

	// Define o Grid como a raiz da aplicação
	if err := i.app.SetRoot(grid, true).EnableMouse(true).Run(); err != nil {
		panic(err)
	}
}

func (i *Terminal) ListenServerResponse(doneCh chan struct{}) {
	for {
		select {
		case <-doneCh:
			return
		case response := <-*i.serverCommandCh:
			// Imprime a resposta recebida
			i.Print(response, i.receivedResponses)
			i.app.Draw()
		}
	}
}

func (i *Terminal) Print(value string, view *tview.TextView) {
	data := time.Now().Format("2006-01-02 15:04:05") + " - " + value + "\n"
	view.Write([]byte(data))
	view.ScrollToEnd()
}

func (i *Terminal) PrintInput(value string) {
	i.mutex.Lock()
	i.app.QueueUpdate(func() {
		i.Print(value, i.sentCommands)
	})
	i.mutex.Unlock()
}

func (i *Terminal) ClearInput() {
	i.mutex.Lock()
	i.app.QueueUpdate(func() {
		i.sentCommands.Clear()
	})
	i.mutex.Unlock()
}

func (i *Terminal) ClearOutput() {
	i.mutex.Lock()
	i.app.QueueUpdate(func() {
		i.receivedResponses.Clear()
	})
	i.mutex.Unlock()
}
