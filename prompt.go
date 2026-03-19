package main

import (
	"bufio"
	"fmt"
	"strings"
)

func prompt(scanner *bufio.Scanner, label string) string {
	fmt.Printf("  %s: ", label)
	scanner.Scan()
	return strings.TrimSpace(scanner.Text())
}

func promptRequired(scanner *bufio.Scanner, label string) string {
	for {
		val := prompt(scanner, label)
		if val != "" {
			return val
		}
		fmt.Printf("    %s cannot be empty.\n", label)
	}
}

func promptWithDefault(scanner *bufio.Scanner, label, current string) string {
	if current != "" {
		fmt.Printf("  %s [%s]: ", label, current)
	} else {
		fmt.Printf("  %s: ", label)
	}
	scanner.Scan()
	input := strings.TrimSpace(scanner.Text())
	if input == "" {
		return current
	}
	return input
}

func promptRequiredWithDefault(scanner *bufio.Scanner, label, current string) string {
	for {
		val := promptWithDefault(scanner, label, current)
		if val != "" {
			return val
		}
		fmt.Printf("    %s cannot be empty.\n", label)
	}
}
