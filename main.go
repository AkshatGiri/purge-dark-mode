package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

func main() {
	// Define command line flags
	dir := flag.String("dir", ".", "Directory to process")
	dryRun := flag.Bool("dry-run", false, "Show what would be changed without making changes")
	flag.Parse()

	// Walk through the directory
	err := filepath.Walk(*dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Skip directories
		if info.IsDir() {
			return nil
		}

		// Process the file
		return processFile(path, *dryRun)
	})

	if err != nil {
		fmt.Fprintf(os.Stderr, "Error walking directory: %v\n", err)
		os.Exit(1)
	}
}

func processFile(path string, dryRun bool) error {
	// Get file info first
	fileInfo, err := os.Stat(path)
	if err != nil {
		return fmt.Errorf("error getting file info for %s: %v", path, err)
	}

	// Read file content
	content, err := ioutil.ReadFile(path)
	if err != nil {
		return fmt.Errorf("error reading file %s: %v", path, err)
	}

	// Convert to string for processing
	originalContent := string(content)
	
	// Create regex to match dark mode classes
	// This matches 'dark:' followed by any characters until a space or any type of quote
	re := regexp.MustCompile(`dark:[^ "'` + "`" + `\n]+`)
	
	// Find all matches to report changes
	matches := re.FindAllString(originalContent, -1)
	if len(matches) == 0 {
		return nil
	}

	// Process the content
	newContent := re.ReplaceAllString(originalContent, "")
	
	// Clean up any double spaces that might have been created
	newContent = strings.ReplaceAll(newContent, "  ", " ")
	
	// Report changes
	fmt.Printf("File: %s\n", path)
	fmt.Printf("Found %d dark mode classes:\n", len(matches))
	for _, match := range matches {
		fmt.Printf("  - %s\n", match)
	}

	// If it's a dry run, don't make changes
	if dryRun {
		fmt.Println("Dry run - no changes made")
		return nil
	}

	// Write the changes back to the file
	err = ioutil.WriteFile(path, []byte(newContent), fileInfo.Mode())
	if err != nil {
		return fmt.Errorf("error writing file %s: %v", path, err)
	}

	fmt.Println("Changes applied successfully")
	return nil
}