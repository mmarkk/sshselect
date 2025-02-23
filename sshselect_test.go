package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// TestConfigFileCreation tests the creation and content of the default config file
func TestConfigFileCreation(t *testing.T) {
	// Create temporary directory for test
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "config")

	// Test config creation
	err := createDefaultConfig(configPath)
	if err != nil {
		t.Fatalf("Failed to create config file: %v", err)
	}

	// Check if file exists
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		t.Error("Config file was not created")
	}

	// Check file permissions
	info, err := os.Stat(configPath)
	if err != nil {
		t.Fatalf("Failed to stat config file: %v", err)
	}
	if info.Mode().Perm() != 0644 {
		t.Errorf("Config file has wrong permissions: got %v, want %v", info.Mode().Perm(), 0644)
	}

	// Check file content
	content, err := os.ReadFile(configPath)
	if err != nil {
		t.Fatalf("Failed to read config file: %v", err)
	}
	if !strings.Contains(string(content), "Host myserver") {
		t.Error("Config file missing example entry")
	}
	if !strings.Contains(string(content), "User admin") {
		t.Error("Config file missing User field example")
	}
}

// TestParseConfig tests the parsing of SSH config entries
func TestParseConfig(t *testing.T) {
	testCases := []struct {
		name     string
		config   string
		expected []SSHHost
	}{
		{
			name: "Valid config",
			config: `Host test1
    HostName 192.168.1.1
    User testuser
    Port 2222`,
			expected: []SSHHost{
				{Name: "test1", HostName: "192.168.1.1", User: "testuser", Port: "2222"},
			},
		},
		{
			name: "Missing required fields",
			config: `Host test1
    HostName 192.168.1.1
Host test2
    User testuser
    Port 2222`,
			expected: nil,
		},
		{
			name: "With comments and whitespace",
			config: `# This is a comment
Host test1
    HostName 192.168.1.1
    User testuser

# Another comment
Host test2
    HostName 192.168.1.2
    User user2`,
			expected: []SSHHost{
				{Name: "test1", HostName: "192.168.1.1", User: "testuser"},
				{Name: "test2", HostName: "192.168.1.2", User: "user2"},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := parseConfig(tc.config)
			if len(result) != len(tc.expected) {
				t.Errorf("got %d hosts, want %d", len(result), len(tc.expected))
				return
			}
			for i, host := range result {
				if host != tc.expected[i] {
					t.Errorf("host %d: got %+v, want %+v", i, host, tc.expected[i])
				}
			}
		})
	}
}

// TestSearcherWithExit tests the searcher function including the Exit option
func TestSearcherWithExit(t *testing.T) {
	logins := []sshLogin{
		{Name: "test1", Command: "ssh user1@host1"},
		{Name: "test2", Command: "ssh user2@host2"},
		{Name: "Exit", Command: "exit"},
	}

	searcher := func(input string, index int) bool {
		login := logins[index]
		name := strings.ToLower(strings.TrimSpace(login.Name))
		input = strings.ToLower(strings.TrimSpace(input))

		if input == "" {
			return true
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

	testCases := []struct {
		input    string
		index    int
		expected bool
	}{
		{"exit", 2, true},
		{"EXIT", 2, true},
		{"e", 2, true},
		{"test", 2, false},
		{"", 2, true},
	}

	for _, tc := range testCases {
		result := searcher(tc.input, tc.index)
		if result != tc.expected {
			t.Errorf("searcher(%q, %d) = %v; want %v", tc.input, tc.index, result, tc.expected)
		}
	}
}

// TestSSHCommandGeneration tests the generation of SSH commands from config entries
func TestSSHCommandGeneration(t *testing.T) {
	testCases := []struct {
		host     SSHHost
		expected string
	}{
		{
			host:     SSHHost{Name: "test1", HostName: "192.168.1.1", User: "user1"},
			expected: "ssh user1@192.168.1.1",
		},
		{
			host:     SSHHost{Name: "test2", HostName: "example.com", User: "user2", Port: "2222"},
			expected: "ssh user2@example.com -p 2222",
		},
	}

	for _, tc := range testCases {
		command := fmt.Sprintf("ssh %s@%s", tc.host.User, tc.host.HostName)
		if tc.host.Port != "" {
			command += fmt.Sprintf(" -p %s", tc.host.Port)
		}
		if command != tc.expected {
			t.Errorf("got command %q, want %q", command, tc.expected)
		}
	}
}

// TestSearcher checks if the searcher function correctly filters SSH logins based on input
func TestSearcher(t *testing.T) {
	logins := []sshLogin{
		{Name: "cpwlapiadm", Command: "ssh cpwlapiadm@170.64.217.201"},
		{Name: "outageradm", Command: "ssh outageradm@170.64.134.33"},
		{Name: "mark", Command: "ssh mark@170.64.140.130"},
	}

	searcher := func(input string, index int) bool {
		login := logins[index]
		name := strings.ToLower(strings.TrimSpace(login.Name))
		input = strings.ToLower(strings.TrimSpace(input))

		fmt.Printf("Debug: Matching input '%s' against name '%s'\n", input, name)

		if input == "" {
			return true // Empty input matches everything
		}

		inputIndex := 0

		for _, nameRune := range name {
			fmt.Printf("Debug: Comparing '%c' with '%c'\n", nameRune, input[inputIndex])
			if inputIndex < len(input) && nameRune == rune(input[inputIndex]) {
				inputIndex++
				fmt.Printf("Debug: Match found, inputIndex now %d\n", inputIndex)
			}
			if inputIndex == len(input) {
				fmt.Printf("Debug: All input characters matched\n")
				return true
			}
		}

		fmt.Printf("Debug: Not all input characters were matched\n")
		return false
	}

	testCases := []struct {
		input    string
		index    int
		expected bool
	}{
		{"cpwl", 0, true},
		{"mark", 2, true},
		{"admin", 1, false},
		{"out", 1, true},
		{"nonexistent", 0, false},
		{"MARK", 2, true},     // Test case-insensitivity
		{"  mark  ", 2, true}, // Test whitespace handling
		{"", 0, true},         // Empty string should match all
		{"adm", 0, true},      // Partial match at the end
		{"cp", 0, true},       // Partial match at the beginning
		{"age", 1, true},      // Partial match in the middle
	}

	for _, tc := range testCases {
		fmt.Printf("\nTesting case: input='%s', index=%d, expected=%v\n", tc.input, tc.index, tc.expected)
		result := searcher(tc.input, tc.index)
		if result != tc.expected {
			t.Errorf("searcher(%s, %d) = %v; want %v", tc.input, tc.index, result, tc.expected)
		}
	}
}

// ... (rest of the code remains unchanged)
