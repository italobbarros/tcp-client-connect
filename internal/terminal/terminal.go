package Interface

import (
	"time"

	"github.com/rivo/tview"
)

type Interface struct {
	serverCommandCh   *chan string
	userCommandCh     *chan string
	sentCommands      *tview.TextView
	receivedResponses *tview.TextView
	app               *tview.Application
}

func NewInterface(serverCommandCh *chan string, userCommandCh *chan string) *Interface {
	return &Interface{
		serverCommandCh: serverCommandCh,
		userCommandCh:   userCommandCh,
	}
}

func (i *Interface) Create() {
	i.app = tview.NewApplication()

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
	form := tview.NewForm().
		AddInputField("Enviar valor", "", 100, nil, nil).
		AddButton("Enviar", nil)
	//AddButton("Sair", func() {
	//	app.Stop()
	//})
	form.SetBorder(true).SetTitle("Enviar dados").SetTitleAlign(tview.AlignLeft)

	form.GetButton(0).SetSelectedFunc(func() {
		command := form.GetFormItemByLabel("Enviar valor").(*tview.InputField).GetText()
		*i.userCommandCh <- command
		i.Print(command, i.sentCommands)
	})
	// Cria um novo Grid e adiciona o Form, os TextViews e um espaço vazio
	grid := tview.NewGrid().
		AddItem(i.sentCommands, 0, 0, 2, 1, 0, 0, false).
		AddItem(i.receivedResponses, 0, 1, 2, 1, 0, 0, false).
		AddItem(form, 2, 0, 1, 2, 0, 0, true)

	// Define o Grid como a raiz da aplicação
	if err := i.app.SetRoot(grid, true).EnableMouse(true).Run(); err != nil {
		panic(err)
	}
}

func (i *Interface) ListenServerResponse() {
	for {
		select {
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
}
