// Package main provides a simple SSH login selector tool.
package main

import (
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"

	"github.com/manifoldco/promptui"
)

type sshLogin struct {
	Name    string
	Command string
}

func main() {
	// SSH login list
	logins := []sshLogin{
		{Name: "cpwlapiadm", Command: "ssh cpwlapiadm@170.64.217.201"},
		{Name: "outageradm", Command: "ssh outageradm@170.64.134.33"},
		{Name: "mark", Command: "ssh mark@170.64.140.130"},
		{Name: "cpwtadm", Command: "ssh cpwtadm@170.64.185.37"},
		{Name: "cpwebadm", Command: "ssh cpwebadm@170.64.132.151"},
	}

	fmt.Println("SSH Login Selector")
	fmt.Println("------------------")

	templates := &promptui.SelectTemplates{
		Label:    "{{ . }}",
		Active:   "\u25B6 {{ .Name | cyan }}",
		Inactive: "  {{ .Name | white }}",
		Selected: "\u2714 {{ .Name | green }}",
		Details:  "Command: {{ .Command }}",
	}

	searcher := func(input string, index int) bool {
		login := logins[index]
		name := strings.ToLower(strings.TrimSpace(login.Name))
		input = strings.ToLower(strings.TrimSpace(input))

		if input == "" {
			return true // Empty input matches everything
		}

		inputRunes := []rune(input)
		nameRunes := []rune(name)
		inputIndex := 0

		for _, nameRune := range nameRunes {
			if inputIndex < len(inputRunes) && nameRune == inputRunes[inputIndex] {
				inputIndex++
			}
			if inputIndex == len(inputRunes) {
				return true
			}
		}

		return false
	}

	prompt := promptui.Select{
		Label:     "Select an SSH login (use arrow keys or type the number):",
		Items:     logins,
		Templates: templates,
		Size:      10,
		Searcher:  searcher,
	}

	index, _, err := prompt.Run()

	if err != nil {
		if err.Error() == "^C" {
			fmt.Println("\nOperation cancelled")
			return
		}
		// Check if the input is a valid number
		if num, err := strconv.Atoi(err.Error()); err == nil && num > 0 && num <= len(logins) {
			index = num - 1
		} else {
			fmt.Printf("Error: %v\n", err)
			os.Exit(1)
		}
	}

	selectedLogin := logins[index]
	fmt.Printf("\nConnecting to: %s\n", selectedLogin.Command)
	sshCommand := strings.Fields(selectedLogin.Command)
	cmd := exec.Command(sshCommand[0], sshCommand[1:]...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err = cmd.Run()
	if err != nil {
		fmt.Println("Error executing SSH command:", err)
		os.Exit(1)
	}
}