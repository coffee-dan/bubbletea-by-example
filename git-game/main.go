package main

import (
	"fmt"
	"os"
	"os/exec"

	tea "github.com/charmbracelet/bubbletea"
)

type model struct {
	longest_streak int
	streak         int
	commit_shas    []string
	committers     []string
	err            error
}

type errMsg struct{ err error }

func (e errMsg) Error() string { return e.err.Error() }

func (mod model) Init() tea.Cmd {
	// mod.commit_shas =
	cmd := exec.Command("git", "log", "--pretty=%H", "-z")
	err := cmd.Run()
	if err != nil {
		mod.err = err
		return nil
	}

	return nil
}

func (mod model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case errMsg:
		mod.err = msg
		return mod, tea.Quit

	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return mod, tea.Quit
		}
	}
	return mod, nil
}

func (mod model) View() string {
	if mod.err != nil {
		return fmt.Sprintf("\nWe had some trouble: %v\n\n", mod.err)
	}

	return "\nHello, World!"
}

func main() {
	if _, err := tea.NewProgram(model{}).Run(); err != nil {
		fmt.Printf("Fuck, %v\n", err)
		os.Exit(1)
	}
}
