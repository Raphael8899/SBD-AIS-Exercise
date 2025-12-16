package main

import (
	"exc9/mapred"
	"fmt"
	"os"
	"strings"
)

// Main function
func main() {
	// 1. Read file contents into memory
	// We read the file "res/meditations.txt" which is located relative to the project root
	content, err := os.ReadFile("res/meditations.txt")
	if err != nil {
		// If the file is missing or cannot be read, we print the error and stop
		fmt.Println("Error reading file:", err)
		return
	}

	// 2. Prepare the input data
	// The MapReduce Run function expects a slice of strings ([]string).
	// Since ReadFile gives us bytes, we convert to string and split by newlines.
	fileContent := string(content)
	inputLines := strings.Split(fileContent, "\n")

	fmt.Println("File read successfully. Starting MapReduce...")

	// 3. Run the MapReduce algorithm
	var mr mapred.MapReduce
	results := mr.Run(inputLines)

	// 4. Print the result to stdout
	// We iterate over the map and print every word and its count.
	fmt.Println("--------------------------------------------------")
	fmt.Println("Word Frequency Results:")
	fmt.Println("--------------------------------------------------")
	
	for word, count := range results {
		fmt.Printf("%s: %d\n", word, count)
	}

	fmt.Println("--------------------------------------------------")
	fmt.Printf("Total unique words found: %d\n", len(results))
}