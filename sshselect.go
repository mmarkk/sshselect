// Package main provides a simple SSH login selector tool.
package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/manifoldco/promptui"
)

type SSHHost struct {
	Name     string
	HostName string
	User     string
	Port     string
}

type sshLogin struct {
	Name    string
	Command string
}

func createDefaultConfig(configPath string) error {
	defaultConfig := `# SSH Login Configuration
# 
# Format: Use standard SSH config format
# Example entries:
#
# Host myserver                    # Nickname for the connection
#     HostName 192.168.1.100      # IP address or hostname
#     User admin                   # SSH username
#     Port 22                     # Optional, defaults to 22
#
# Host aws-instance
#     HostName ec2-1-2-3-4.compute-1.amazonaws.com
#     User ubuntu
#     Port 2222
#
# Add your SSH connections below:

`
	// Ensure directory exists
	configDir := filepath.Dir(configPath)
	fmt.Printf("Creating config directory: %s\n", configDir)
	if err := os.MkdirAll(configDir, 0750); err != nil {
		return fmt.Errorf("failed to create config directory: %v", err)
	}

	// Create config file with default content
	fmt.Printf("Writing default config to: %s\n", configPath)
	if err := os.WriteFile(configPath, []byte(defaultConfig), 0600); err != nil {
		return fmt.Errorf("failed to write config file: %v", err)
	}
	fmt.Printf("Successfully created config file\n")
	return nil
}

func parseConfig(content string) []SSHHost {
	var hosts []SSHHost
	var currentHost *SSHHost

	scanner := bufio.NewScanner(strings.NewReader(content))
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())

		// Skip empty lines and comments
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		// Check if this is a Host line
		if strings.HasPrefix(strings.ToLower(line), "host ") {
			// Add previous host if valid
			if currentHost != nil && currentHost.User != "" && currentHost.HostName != "" {
				hosts = append(hosts, *currentHost)
			} else if currentHost != nil {
				fmt.Printf("\nWarning: Skipping '%s' - missing required fields (User and/or HostName)\n", currentHost.Name)
			}
			hostName := strings.TrimSpace(strings.TrimPrefix(line, "Host "))
			currentHost = &SSHHost{Name: hostName}
			continue
		}

		// Parse other config lines
		if currentHost != nil {
			parts := strings.Fields(line)
			if len(parts) >= 2 {
				key := strings.ToLower(parts[0])
				value := strings.Join(parts[1:], " ")

				switch key {
				case "hostname":
					currentHost.HostName = value
				case "user":
					currentHost.User = value
				case "port":
					currentHost.Port = value
				}
			}
		}
	}

	// Add the last host if valid
	if currentHost != nil && currentHost.User != "" && currentHost.HostName != "" {
		hosts = append(hosts, *currentHost)
	} else if currentHost != nil {
		fmt.Printf("\nWarning: Skipping '%s' - missing required fields (User and/or HostName)\n", currentHost.Name)
	}

	if len(hosts) == 0 {
		return nil
	}

	return hosts
}

func loadConfig() ([]sshLogin, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return nil, fmt.Errorf("failed to get home directory: %v", err)
	}

	// Construct and validate config path
	configPath := filepath.Join(homeDir, ".config", "sshselect", "config")
	expectedPrefix := filepath.Join(homeDir, ".config", "sshselect")
	if !strings.HasPrefix(filepath.Clean(configPath), expectedPrefix) {
		return nil, fmt.Errorf("config file must be in %s directory", expectedPrefix)
	}

	fmt.Printf("Checking config at: %s\n", configPath)

	// Check if config exists
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		fmt.Printf("Config file not found, creating...\n")
		if err := createDefaultConfig(configPath); err != nil {
			return nil, fmt.Errorf("failed to create default config: %v", err)
		}

		// Verify the file was created
		if _, err := os.Stat(configPath); err != nil {
			return nil, fmt.Errorf("failed to verify config file creation: %v", err)
		}

		fmt.Printf("\nCreated default config file at: %s\nPlease add your SSH connections and run the program again.\n", configPath)
		return nil, fmt.Errorf("new config file created")
	} else {
		fmt.Printf("Using existing config file: %s\n", configPath)
	}

	// Read config file
	content, err := os.ReadFile(configPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %v", err)
	}

	// Parse hosts from config
	hosts := parseConfig(string(content))

	// Convert to sshLogin format
	var logins []sshLogin
	for _, host := range hosts {
		command := fmt.Sprintf("ssh %s@%s", host.User, host.HostName)
		if host.Port != "" {
			command += fmt.Sprintf(" -p %s", host.Port)
		}
		logins = append(logins, sshLogin{
			Name:    host.Name,
			Command: command,
		})
	}

	if len(logins) == 0 {
		fmt.Printf("No valid SSH connections found in config file: %s\n", configPath)
		fmt.Printf("Each entry must have:\n")
		fmt.Printf("  Host nickname\n")
		fmt.Printf("      HostName <ip-or-hostname>\n")
		fmt.Printf("      User <username>\n")
		fmt.Printf("      Port <port>    # Optional\n")
		fmt.Printf("\nPlease update your connections and run the program again.\n")
		os.Exit(0)
	}

	return logins, nil
}

func main() {
	fmt.Println("SSH Login Selector")
	fmt.Println("------------------")

	homeDir, err := os.UserHomeDir()
	if err != nil {
		fmt.Printf("Error getting home directory: %v\n", err)
		os.Exit(1)
	}

	// Construct and validate config path
	configPath := filepath.Join(homeDir, ".config", "sshselect", "config")
	expectedPrefix := filepath.Join(homeDir, ".config", "sshselect")
	if !strings.HasPrefix(filepath.Clean(configPath), expectedPrefix) {
		fmt.Printf("Error: config file must be in %s directory\n", expectedPrefix)
		os.Exit(1)
	}

	// Check if config exists
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		fmt.Printf("Config file not found, creating at: %s\n", configPath)
		if err := createDefaultConfig(configPath); err != nil {
			fmt.Printf("Error creating config: %v\n", err)
			os.Exit(1)
		}
		fmt.Printf("\nCreated default config file at: %s\nPlease add your SSH connections and run the program again.\n", configPath)
		return
	}

	logins, err := loadConfig()
	if err != nil {
		fmt.Printf("Error loading config: %v\n", err)
		os.Exit(1)
	}

	// Add exit option
	logins = append(logins, sshLogin{
		Name:    "Exit",
		Command: "exit",
	})

	templates := &promptui.SelectTemplates{
		Label:    "{{ . }}",
		Active:   `{{ if eq .Name "Exit" }}` + "\u25B6 {{ .Name | red }}" + `{{ else }}` + "\u25B6 {{ .Name | cyan }}" + `{{ end }}`,
		Inactive: `{{ if eq .Name "Exit" }}` + "  {{ .Name | red }}" + `{{ else }}` + "  {{ .Name | white }}" + `{{ end }}`,
		Selected: "\u2714 {{ .Name | green }}",
		Details:  `{{ if ne .Name "Exit" }}Command: {{ .Command }}{{ end }}`,
	}

	searcher := func(input string, index int) bool {
		login := logins[index]
		name := strings.ToLower(strings.TrimSpace(login.Name))
		input = strings.ToLower(strings.TrimSpace(input))

		if input == "" {
			return true // Empty input matches everything
		}

		// Use string iteration directly
		inputIdx := 0
		for _, r := range input {
			found := false
			for i := 0; i < len(name); i++ {
				if rune(name[i]) == r {
					found = true
					break
				}
			}
			if !found {
				return false
			}
			inputIdx++
		}
		return true
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
	if selectedLogin.Command == "exit" {
		fmt.Println("\nExiting...")
		return
	}

	fmt.Printf("\nConnecting to: %s\n", selectedLogin.Command)
	sshCommand := strings.Fields(selectedLogin.Command)

	// Validate and extract SSH arguments
	if len(sshCommand) < 2 {
		fmt.Println("Error: Invalid SSH command format")
		os.Exit(1)
	}

	// Expect format: ssh user@host [-p port]
	var sshArgs []string
	userHost := sshCommand[1]
	if !strings.Contains(userHost, "@") {
		fmt.Println("Error: Invalid SSH command format - missing user@host")
		os.Exit(1)
	}

	// Split user@host and validate
	parts := strings.Split(userHost, "@")
	if len(parts) != 2 || parts[0] == "" || parts[1] == "" {
		fmt.Println("Error: Invalid SSH command format - malformed user@host")
		os.Exit(1)
	}

	// Add validated user@host
	sshArgs = append(sshArgs, userHost)

	// Handle optional port argument
	if len(sshCommand) > 2 {
		if len(sshCommand) != 4 || sshCommand[2] != "-p" {
			fmt.Println("Error: Invalid SSH command format - unexpected arguments")
			os.Exit(1)
		}
		// Validate port is numeric
		if _, err := strconv.Atoi(sshCommand[3]); err != nil {
			fmt.Println("Error: Invalid SSH command format - port must be numeric")
			os.Exit(1)
		}
		sshArgs = append(sshArgs, "-p", sshCommand[3])
	}

	// Use explicit ssh command with validated arguments
	cmd := exec.Command("ssh", sshArgs...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err = cmd.Run()
	if err != nil {
		fmt.Println("Error executing SSH command:", err)
		os.Exit(1)
	}
}
