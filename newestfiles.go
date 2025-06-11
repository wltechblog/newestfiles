package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"
)

type FileInfo struct {
	Path    string
	ModTime time.Time
}

func main() {
	// Define command line flags
	jsonOutput := flag.Bool("j", false, "output in JSON format")
	flag.Parse()

	// Get suffix arguments from command line (after flags)
	suffixes := flag.Args()

	// If no suffixes provided, show usage
	if len(suffixes) == 0 {
		fmt.Println("Usage: newestfiles [-j] <suffix1> [suffix2] ...")
		fmt.Println("  -j    output in JSON format (default: plain text)")
		fmt.Println("Example: newestfiles .txt .go .md")
		fmt.Println("Example: newestfiles -j .txt .go .md")
		return
	}

	// Normalize suffixes to ensure they start with a dot
	for i, suffix := range suffixes {
		if !strings.HasPrefix(suffix, ".") {
			suffixes[i] = "." + suffix
		}
	}

	var files []FileInfo

	// Walk through current directory and subdirectories
	err := filepath.Walk(".", func(path string, info os.FileInfo, err error) error {
		if err != nil {
			fmt.Printf("Error accessing %s: %v\n", path, err)
			return nil // Continue walking despite errors
		}

		// Skip directories
		if info.IsDir() {
			return nil
		}

		// Check if file has one of the target suffixes
		for _, suffix := range suffixes {
			if strings.HasSuffix(strings.ToLower(info.Name()), strings.ToLower(suffix)) {
				files = append(files, FileInfo{
					Path:    path,
					ModTime: info.ModTime(),
				})
				break // Found a match, no need to check other suffixes
			}
		}

		return nil
	})

	if err != nil {
		fmt.Printf("Error walking directory: %v\n", err)
		return
	}

	if len(files) == 0 {
		fmt.Println("No files found with the specified suffixes.")
		return
	}

	// Sort files by modification time in descending order (newest first)
	sort.Slice(files, func(i, j int) bool {
		return files[i].ModTime.After(files[j].ModTime)
	})

	// Output the sorted list
	if *jsonOutput {
		// JSON output
		var fns []string
		for _, file := range files {
			fns = append(fns, file.Path)
		}
		out, _ := json.Marshal(&fns)
		fmt.Printf("%s", out)
	} else {
		// Plain text output (default)
		for _, file := range files {
			fmt.Println(file.Path)
		}
	}
}
