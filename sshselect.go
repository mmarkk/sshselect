// Package main provides a simple SSH login selector tool.
package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

// main is the entry point of the program. It presents SSH logins to the user,
// and executes the selected SSH command.
func main() {
	// SSH login list
	logins := []string{
		"ssh cpwtadm@170.64.185.37",
		"ssh cpwebadm@170.64.132.151",
		"ssh mark@170.64.140.130",
		"ssh outageradm@170.64.134.33",
		"ssh cpwlapiadm@170.64.217.201",
	}

	fmt.Println("SSH Login Selector")
	fmt.Println("------------------")
	// Present options to the user
	for i, login := range logins {
		fmt.Printf("%d. %s\n", i+1, login)
	}

	// Get user selection
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("\nSelect an SSH login (enter the number): ")
	input, err := reader.ReadString('\n')
	if err != nil {
		fmt.Println("Error reading input:", err)
		os.Exit(1)
	}

	// Parse user input
	choice := strings.TrimSpace(input)
	index := 0
	_, err = fmt.Sscanf(choice, "%d", &index)
	if err != nil || index < 1 || index > len(logins) {
		fmt.Println("Invalid selection")
		os.Exit(1)
	}

	// Execute SSH command
	selectedLogin := logins[index-1]
	fmt.Printf("\nConnecting to: %s\n", selectedLogin)
	sshCommand := strings.Fields(selectedLogin)
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