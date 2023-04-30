package main

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"

	git "github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/object"
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
	repo, err := git.PlainClone('.', false)

	cIter, err := repo.Log(&git.LogOptions{})

	// ... just iterates over the commits, printing it
	err = cIter.ForEach(func(c *object.Commit) error {
		fmt.Println(c)

		return nil
	})
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
