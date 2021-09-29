package config

import (
	"bufio"
	"os"
	"strings"
)

func GetNodesList(nodesListFile string) []string {
	f, _ := os.Open(nodesListFile)
	// Create new Scanner.
	scanner := bufio.NewScanner(f)
	result := []string{}
	// Use Scan.
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		// Append line to result.
		if line != "" {
			result = append(result, line)
		}
	}
	return result
}
