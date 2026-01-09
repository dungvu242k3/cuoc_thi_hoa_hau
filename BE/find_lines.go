package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func main() {
	file, err := os.Open("internal/adapter/graph/generated/generated.go")
	if err != nil {
		panic(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	line := 1
	for scanner.Scan() {
		text := scanner.Text()
		if strings.Contains(text, "parsedSchema") {
			fmt.Printf("Line %d: %s\n", line, text)
		}
		if strings.Contains(text, "sources") && strings.Contains(text, "ast.Source") {
			fmt.Printf("Line %d: %s\n", line, text)
		}
		line++
	}
}
