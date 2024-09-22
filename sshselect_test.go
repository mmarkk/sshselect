package main

import (
	"fmt"
	"strings"
	"testing"
)

// ... (previous code remains unchanged)

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
		{"MARK", 2, true},  // Test case-insensitivity
		{"  mark  ", 2, true},  // Test whitespace handling
		{"", 0, true},  // Empty string should match all
		{"adm", 0, true},  // Partial match at the end
		{"cp", 0, true},  // Partial match at the beginning
		{"age", 1, true},  // Partial match in the middle
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