package main

import (
	"fmt"
	"os"

	"github.com/go-git/go-git/plumbing/object"
	git "github.com/go-git/go-git/v5"
)

func checkIfError(err error) {
	if err == nil {
		return
	}

	fmt.Printf("\x1b[31;1m%s\x1b[0m\n", fmt.Sprintf("error: %s", err))
	os.Exit(1)
}

func main() {
	repo, err := git.PlainOpen("..")
	checkIfError(err)

	cIter, err := repo.Log(&git.LogOptions{})
	checkIfError(err)

	err = cIter.ForEach(func(c *object.Commit) error {
		fmt.Println(c)

		return nil
	})
	checkIfError(err)
}
