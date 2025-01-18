package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
)

func main() {
	dir := flag.String("dir", ".", "Directory to process")
	dryRun := flag.Bool("dry-run", false, "Show what would be changed without making changes")
	flag.Parse()

	err := filepath.Walk(*dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() {
			return nil
		}

		return processFile(path, *dryRun)
	})

	if err != nil {
		fmt.Fprintf(os.Stderr, "Error walking directory: %v\n", err)
		os.Exit(1)
	}
}

func processFile(path string, dryRun bool) error {
	fileInfo, err := os.Stat(path)
	if err != nil {
		return fmt.Errorf("error getting file info for %s: %v", path, err)
	}

	content, err := ioutil.ReadFile(path)
	if err != nil {
		return fmt.Errorf("error reading file %s: %v", path, err)
	}

	str := string(content)
	result := RemoveDarkClass(str)

	if dryRun {
		fmt.Println("Dry run - no changes made")
		return nil
	}

	err = ioutil.WriteFile(path, []byte(result), fileInfo.Mode())
	if err != nil {
		return fmt.Errorf("error writing file %s: %v", path, err)
	}

	fmt.Println("Changes applied successfully")
	return nil
}

func RemoveDarkClass(str string) string {
	result := []byte{}

	for i := 0; i < len(str); i++ {
		start := str[i]
		if start != 'd' {
			result = append(result, start)
			continue
		}

		isDark := str[i:i+5] == "dark:"

		if !isDark {
			result = append(result, start)
			continue
		}

		for j := i + 5; j < len(str); j++ {
			current := str[j]
			// if we find a space, we'll include that in the skip
			if current == ' ' || current == '\n' || current == '\t' || current == '\r' {
				// keep looping until we find a non-space character
				for k := j; k < len(str); k++ {
					if str[k] != ' ' && str[k] != '\n' && str[k] != '\t' && str[k] != '\r' {
						j = k - 1
						break
					}
				}
				i = j
				break
				// if we find a quote or backtick, we will include everything up to that point.
			} else if current == '"' || current == '\'' || current == '`' {
				// check if the character before the start is a space, we need to remove that too
				if str[i-1] == ' ' {
					result = result[:len(result)-1]
				}
				i = j - 1
				break
			}
		}
	}

	return string(result)
}
