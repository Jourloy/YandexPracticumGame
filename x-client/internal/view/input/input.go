package input

// A simple program demonstrating the text input component from the Bubbles
// component library.

import (
	"fmt"
	"log"
	"os"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/jourloy/x-client/internal/config/env"
)

func Main() {
	p := tea.NewProgram(initialModel())
	if _, err := p.Run(); err != nil {
		log.Fatal(err)
	}
}

type (
	errMsg error
)

type model struct {
	textInput textinput.Model
	err       error
}

func initialModel() model {
	ti := textinput.New()
	ti.Placeholder = `COOL AND LONG API KEY`
	ti.Focus()
	ti.CharLimit = 156
	ti.Width = 20

	return model{
		textInput: ti,
		err:       nil,
	}
}

func (m model) Init() tea.Cmd {
	return textinput.Blink
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyEnter:
			if _, err := os.Stat(`./.env`); os.IsNotExist(err) {
			} else {
				if err := os.Truncate(`./.env`, 0); err != nil {
					return model{}, tea.Quit
				}
			}

			file, err := os.OpenFile(`./.env`, os.O_RDWR|os.O_CREATE, 0666)
			if err != nil {
				return model{}, tea.Quit
			}
			defer file.Close()

			data := []byte(`API_KEY=` + m.textInput.Value())
			file.Write(data)

			env.API = m.textInput.Value()

			return m, nil
		case tea.KeyEsc, tea.KeyCtrlC:
			return m, tea.Quit
		}

	case errMsg:
		return m, tea.Quit
	}

	m.textInput, cmd = m.textInput.Update(msg)
	return m, cmd
}

func (m model) View() string {
	return fmt.Sprintf(
		"Enter your API key:\n\n%s\n\n%s",
		m.textInput.View(),
		"(esc to quit)",
	) + "\n"
}
