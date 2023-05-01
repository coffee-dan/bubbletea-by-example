package main

import (
	"bytes"
	"fmt"
	"log"
	"math/rand"
	"os/exec"
	"strings"
	"time"
)

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

func getCommitShas() []string {
	cmd := exec.Command("git", "log", "--no-merges", "--pretty=%H", "-z")
	output := getStringOutput(*cmd)
	split_output := strings.Split(output, "\x00")
	return shuffle(split_output)
}

func getCommitAuthor(sha string) string {
	cmd := exec.Command("git", "show", sha, "--pretty=%aN", "--no-patch")
	return getStringOutput(*cmd)
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

func getCommitters() []string {
	cmd := exec.Command("git", "log", "--no-merges", "--pretty=%aN")
	output := getStringOutput(*cmd)
	splitOutput := strings.Split(output, "\n")
	return uniq(splitOutput)
}

func main() {
	commitShas := getCommitShas()
	fmt.Println(commitShas)
	fmt.Println()
	fmt.Println(getCommitters())

	sha := commitShas[len(commitShas)-1]
	commitShas = commitShas[:len(commitShas)-1]
	fmt.Println(sha)

	fmt.Println(getCommitAuthor(sha))

}
