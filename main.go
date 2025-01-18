package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	ignore "github.com/sabhiram/go-gitignore"
)

func main() {

	// Parse command line flags
	dir := flag.String("dir", ".", "Directory to process")
	dryRun := flag.Bool("dry-run", false, "Show what would be changed without making changes")
	logLevelFlag := flag.String("log-level", "info", "Log level (debug, info, warn, error)")

	flag.Parse()

	var logLevel zerolog.Level

	switch *logLevelFlag {
	case "debug":
		logLevel = zerolog.DebugLevel
	case "info":
		logLevel = zerolog.InfoLevel
	case "warn":
		logLevel = zerolog.WarnLevel
	case "error":
		logLevel = zerolog.ErrorLevel
	default:
		log.Error().Msg("Invalid log level")
		os.Exit(1)
	}


	// Setup logging
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	zerolog.SetGlobalLevel(logLevel)
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})

	// Check if .gitignore file exists in the directory
	gitignorePath := filepath.Join(*dir, ".gitignore")

	var ignoreParser *ignore.GitIgnore

	if _, err := os.Stat(gitignorePath); os.IsNotExist(err) {
		log.Warn().Msg("No .gitignore file found in the directory. It's recommended to have a .gitignore file to exclude node_modules and other directories from being processed.")
		ignoreParser = ignore.CompileIgnoreLines()
	} else { 
		ignoreParser, err = ignore.CompileIgnoreFile(gitignorePath)
		if err != nil {
			log.Error().Err(err).Msg("Error compiling .gitignore file")
			os.Exit(1)
		}
	}

	processedFilesCount := 0

	// walk the directory
	err := filepath.Walk(*dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// explicitly skipping .gitignore file
		if !info.IsDir() && info.Name() == ".gitignore" {
			log.Debug().Str("path", path).Msg("Ignoring file")
			return nil
		}

		// explicitly skipping .git folder
		if info.IsDir() && info.Name() == ".git" { 
			log.Debug().Str("path", path).Msg("Ignoring directory")
			return filepath.SkipDir
		}

		// skipping directories that are gitignored
		if info.IsDir() && ignoreParser.MatchesPath(path) { 
			log.Debug().Str("path", path).Msg("Ignoring directory")
			return filepath.SkipDir
		}

		// skipping files that are gitignored
		if(ignoreParser.MatchesPath(path)) {
			log.Debug().Str("path", path).Msg("Ignoring file")
			return nil
		}

		// not processing the folder itself
		if info.IsDir() {
			// just a sanity warning to warn user about non gitignored directories
			if info.Name() == "node_modules" { 
				log.Warn().Msg("Found node_modules while walking directory. Please make sure it's part of the .gitignore. Processing files in node_modules can cause issues. If there other directories you'd like to ignore, please add them to the .gitignore file.")
			}
			return nil
		}

		// process all other type of files
		err = processFile(path, *dryRun)
		processedFilesCount += 1
		return err
	})

	if err != nil {
		log.Error().Err(err).Msg("Error walking directory: " + *dir)
		os.Exit(1)
	}

	if *dryRun {
		log.Info().Msg("Dry run - no changes made")
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

	if(len(str) < 5) { 
		return str, found
	}

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
