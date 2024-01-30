package Interface

import (
	"time"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type Interface struct {
	serverCommandCh   *chan string
	userCommandCh     *chan string
	sentCommands      *tview.TextView
	receivedResponses *tview.TextView
	app               *tview.Application
	data              *tview.Form
}

func NewInterface(serverCommandCh *chan string, userCommandCh *chan string) *Interface {
	return &Interface{
		serverCommandCh: serverCommandCh,
		userCommandCh:   userCommandCh,
	}
}

func (i *Interface) Create(doneCh chan struct{}) {
	i.app = tview.NewApplication()
	go func() {
		for {
			select {
			case <-doneCh:
				i.app.Stop()
			}
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

	//form.GetButton(0).SetSelectedFunc(func() {
	//	command := form.GetFormItemByLabel("Dados").(*tview.InputField).GetText()
	//	if len(command) > 0 {
	//		*i.userCommandCh <- command
	//		i.Print(command, i.sentCommands)
	//		form.GetFormItemByLabel("Dados").(*tview.InputField).SetText("")
	//		i.app.SetFocus(form.GetFormItemByLabel("Dados"))
	//	}
	//})

	i.data.GetFormItemByLabel("Dados").(*tview.InputField).
		SetDoneFunc(func(key tcell.Key) {
			if key == tcell.KeyEnter {
				// Chama a função de envio ao pressionar "Enter"
				command := i.data.GetFormItemByLabel("Dados").(*tview.InputField).GetText()
				if len(command) > 0 {
					*i.userCommandCh <- command
					i.Print(command, i.sentCommands)
					i.data.GetFormItemByLabel("Dados").(*tview.InputField).SetText("")
					i.app.SetFocus(i.data.GetFormItemByLabel("Dados"))

				}
			}
		})
	// Cria um novo Grid e adiciona o Form, os TextViews e um espaço vazio
	grid := tview.NewGrid().
		AddItem(i.sentCommands, 0, 0, 4, 1, 0, 0, false).
		AddItem(i.receivedResponses, 0, 1, 4, 1, 0, 0, false).
		AddItem(i.data, 4, 0, 1, 2, 0, 0, true)

	// Define o Grid como a raiz da aplicação
	if err := i.app.SetRoot(grid, true).EnableMouse(true).Run(); err != nil {
		panic(err)
	}
}

func (i *Interface) ListenServerResponse(doneCh chan struct{}) {
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

func (i *Interface) Print(value string, view *tview.TextView) {
	data := time.Now().Format("2006-01-02 15:04:05") + " - " + value + "\n"
	view.Write([]byte(data))
	view.ScrollToEnd()
}
