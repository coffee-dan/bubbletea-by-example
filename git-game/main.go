package main

import (
	"bytes"
	"fmt"
	"log"
	"math/rand"
	"os"
	"os/exec"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
)

const header = `
----------------------------------------------------------
                      THE GIT GAME
----------------------------------------------------------
Welcome! The goal of the git game is to guess committers
based on their commit messages.


`

type model struct {
	longest_streak int
	streak         int
	commitShas     []string
	committers     []string
	currentCommit  string
	userReady      bool
	choices        []string
	cursor         int
	err            error
}

type errMsg struct{ err error }

func (e errMsg) Error() string { return e.err.Error() }

func getStringOutput(cmd exec.Cmd) string {
	var outBuffer, errBuffer bytes.Buffer

	cmd.Stdout = &outBuffer
	cmd.Stderr = &errBuffer

	err := cmd.Run()
	if err != nil {
		log.Fatalf("command failed with error: %v, stderr: %s", err, errBuffer.Bytes())
	}
	return string(outBuffer.Bytes())
}

func uniq(baseSlice []string) []string {
	uniqueSlice := make([]string, 0, len(baseSlice))
	seenMap := make(map[string]bool)
	for _, str := range baseSlice {
		if !seenMap[str] {
			uniqueSlice = append(uniqueSlice, str)
			seenMap[str] = true
		}
	}
	return uniqueSlice
}

func shuffle(slice []string) []string {
	rand.Seed(time.Now().UnixNano())

	for i := range slice {
		newIdx := rand.Intn(i + 1)
		slice[i], slice[newIdx] = slice[newIdx], slice[i]
	}

	return slice
}

func compactStrip(slice []string) []string {
	filteredSlice := make([]string, 0, len(slice))
	for _, str := range slice {
		cleanStr := strings.TrimSpace(str)
		if len(cleanStr) > 1 {
			filteredSlice = append(filteredSlice, cleanStr)
		}
	}
	return filteredSlice
}

func getCommitShas() []string {
	cmd := exec.Command("git", "log", "--no-merges", "--pretty=%H", "-z")
	output := getStringOutput(*cmd)
	split_output := strings.Split(output, "\x00")
	return shuffle(split_output)
}

func getCommitters() []string {
	cmd := exec.Command("git", "log", "--no-merges", "--pretty=%aN")
	output := getStringOutput(*cmd)
	return compactStrip(uniq(strings.Split(output, "\n")))
}

func getCommitAuthor(sha string) string {
	cmd := exec.Command("git", "show", sha, "--pretty=%aN", "--no-patch")
	return strings.TrimSpace(getStringOutput(*cmd))
}

func getCommitPreview(sha string) string {
	cmd := exec.Command("git", "show", sha, "--pretty=(%ar)%n%B", "--shortstat")
	return getStringOutput(*cmd)
}

func getNextCommit(mod model) string {
	nextCommit := mod.commitShas[len(mod.commitShas)-1]
	mod.commitShas = mod.commitShas[:len(mod.commitShas)-1]
	return nextCommit
}

func getChoices(mod model) []string {
	author := getCommitAuthor(mod.currentCommit)
	choices := []string{author}
	numChoices := 4
	for _, name := range shuffle(mod.committers) {
		if len(choices) >= numChoices {
			break
		}
		if name != author {
			choices = append(choices, name)
		}
	}
	return choices
}

func initialModel() model {
	return model{
		commitShas: getCommitShas(),
		committers: getCommitters(),
		userReady:  false,
	}
}

func (mod model) Init() tea.Cmd {
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

		case "up", "k":
			if mod.cursor > 0 {
				mod.cursor--
			}

		case "down", "j":
			if mod.cursor < len(mod.choices)-1 {
				mod.cursor++
			}

		case "enter":
			if !mod.userReady {
				mod.currentCommit = getNextCommit(mod)
				mod.userReady = true
				mod.choices = getChoices(mod)
			}
		}
	}
	return mod, nil
}

func (mod model) View() string {
	if mod.err != nil {
		return fmt.Sprintf("\nWe had some trouble: %v\n\n", mod.err)
	}

	str := "Git Game\n\n"

	if !mod.userReady {
		str = fmt.Sprintf("You're playing in a repo with %d commits and %d\n", len(mod.commitShas), len(mod.committers))
		str += fmt.Sprintf("distinct committer(s).\n\n")
		// str += fmt.Sprintf("%v\n", mod.committers)
		for _, name := range mod.committers {
			str += fmt.Sprintf("%s\n", name)
		}
		str += "\nReady? PRESS ENTER TO START PLAYING (q to quit)"
		return str
	}
	if len(mod.commitShas) > 0 {
		// str += fmt.Sprintf("%v", mod.commitShas)
		// str += fmt.Sprintf(mod.currentCommit)
		str += fmt.Sprintf(getCommitPreview(mod.currentCommit))
		str += "\n\n"
		for idx, choice := range mod.choices {
			cursor := " "
			if mod.cursor == idx {
				cursor = ">"
			}

			str += fmt.Sprintf("%s [%d] %s\n", cursor, idx+1, choice)
		}
	}

	str += "\nPress q to quit.\n"
	return str
}

func main() {
	prog := tea.NewProgram(initialModel())
	if _, err := prog.Run(); err != nil {
		fmt.Printf("Fuck, %v\n", err)
		os.Exit(1)
	}
}
