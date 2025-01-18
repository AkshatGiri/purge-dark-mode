package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func main() {
	// Parse command line flags
	dir := flag.String("dir", ".", "Directory to process")
	dryRun := flag.Bool("dry-run", false, "Show what would be changed without making changes")
	flag.Parse()

	// Setup logging
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	zerolog.SetGlobalLevel(zerolog.DebugLevel)
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})

	processedFilesCount := 0

	// walk the directory
	err := filepath.Walk(*dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// TODO: check if we need to process the file.
		if info.IsDir() {
			return nil
		}

		// process the file.
		err = processFile(path, *dryRun)
		processedFilesCount += 1
		return err
	})

	if err != nil {
		log.Error().Err(err).Msg("Error walking directory: " + *dir)
		os.Exit(1)
	}

	log.Info().Int("files_processed", processedFilesCount).Msg("Processing complete")
}

func processFile(path string, dryRun bool) error {
	fileInfo, err := os.Stat(path)
	log.Info().Str("path", path).Msg("Processing file")
	if err != nil {
		log.Error().
        Err(err).
        Str("path", path).
        Msg("error getting file info")
		return fmt.Errorf("error getting file info for %s: %v", path, err)
	}

	content, err := ioutil.ReadFile(path)
	if err != nil {
		log.Error().
        Err(err).
        Str("path", path).
        Msg("error reading file contents")
		return fmt.Errorf("error reading file %s: %v", path, err)
	}

	str := string(content)
	result, found := RemoveDarkClasses(str)

	log.Debug().Strs("found", found).Int("count", len(found)).Msg("Found dark classes")

	if dryRun {
		log.Info().Msg("Dry run - no changes made")
		return nil
	}

	err = ioutil.WriteFile(path, []byte(result), fileInfo.Mode())
	if err != nil {
		log.Error().
        Err(err).
        Str("path", path).
        Msg("error writing file")
		return fmt.Errorf("error writing file %s: %v", path, err)
	}

	log.Info().Str("path", path).Msg("File updated")
	return nil
}

func RemoveDarkClasses(str string) (string, []string) {
	result := []byte{}
	found := []string{}

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
				found = append(found, str[i:j+1])
				i = j
				break
				// if we find a quote or backtick, we will include everything up to that point.
			} else if current == '"' || current == '\'' || current == '`' {
				// check if the character before the start is a space, we need to remove that too
				if str[i-1] == ' ' {
					result = result[:len(result)-1]
				}
				found = append(found, str[i:j])
				i = j - 1
				break
			}
		}
	}

	return string(result), found
}
