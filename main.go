package main

import (
	"log"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/gothew/cuptui/code"
	"github.com/gothew/cuptui/filetree"
)

type sessionState int

const (
	idleState sessionState = iota
	showCodeState
)

type Bubble struct {
	filetree filetree.Bubble
	code     code.Bubble
	state    sessionState
}

func (b *Bubble) resetViewports() {
	b.code.Viewport.GotoTop()
}

func (b *Bubble) openFile() []tea.Cmd {
	var cmds []tea.Cmd

	selectedFile := b.filetree.GetSelectedItem()

	if !selectedFile.IsDirectory() {
		b.resetViewports()

		b.state = showCodeState
		readFileCmd := b.code.SetFileName(selectedFile.FileName())
		cmds = append(cmds, readFileCmd)
	}

	return cmds
}

func New() Bubble {
	filetreeModel := filetree.New(
		true,
		"",
		"",
		lipgloss.AdaptiveColor{Light: "#000000", Dark: "63"},
		lipgloss.AdaptiveColor{Light: "#000000", Dark: "63"},
		lipgloss.AdaptiveColor{Light: "63", Dark: "63"},
		lipgloss.AdaptiveColor{Light: "#ffffff", Dark: "#ffffff"},
	)

	codeModel := code.New(true, true, lipgloss.AdaptiveColor{Light: "#000000", Dark: "ffffff"})

	return Bubble{
		filetree: filetreeModel,
		code:     codeModel,
	}
}

func (b Bubble) Init() tea.Cmd {
	return b.filetree.Init()
}

func (b Bubble) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var (
		cmd  tea.Cmd
		cmds []tea.Cmd
	)

	b.filetree, cmd = b.filetree.Update(msg)
	cmds = append(cmds, cmd)

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		halfSize := msg.Width / 2
		b.filetree.SetSize(halfSize, msg.Height)
		b.code.SetSize(halfSize, msg.Height)
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c":
			cmds = append(cmds, tea.Quit)
		case "p":
			cmds = append(cmds, tea.Batch(b.openFile()...))
		}
	}

	b.code, cmd = b.code.Update(msg)
	cmds = append(cmds, cmd)

	return b, tea.Batch(cmds...)
}

func (b Bubble) View() string {
	leftBox := b.filetree.View()
	rightBox := b.code.View()

  switch b.state {
  case idleState:
    rightBox = b.code.View()
  case showCodeState:
    rightBox = b.code.View()
  }

	return lipgloss.JoinVertical(lipgloss.Top,
		lipgloss.JoinHorizontal(lipgloss.Top, leftBox, rightBox),
	)
}

func main() {
	b := New()
	p := tea.NewProgram(b, tea.WithAltScreen())

	if err := p.Start(); err != nil {
		log.Fatal(err)
	}
}
