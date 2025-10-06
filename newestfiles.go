package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"
)

type FileInfo struct {
	Path    string
	ModTime time.Time
	Size    int64
}

func main() {
	// Define command line flags
	jsonOutput := flag.Bool("j", false, "output in JSON format")
	oldest := flag.Bool("o", false, "Sort oldest to newest")
	largest := flag.Bool("l", false, "Sort by largest files first")
	smallest := flag.Bool("s", false, "Sort by smallest files first")
	flag.Parse()

	// Get suffix arguments from command line (after flags)
	suffixes := flag.Args()

	// Check for conflicting sort flags
	sortFlags := 0
	if *oldest {
		sortFlags++
	}
	if *largest {
		sortFlags++
	}
	if *smallest {
		sortFlags++
	}
	if sortFlags > 1 {
		fmt.Println("Error: Only one sort option can be specified at a time")
		return
	}

	// Normalize suffixes to ensure they start with a dot (if any suffixes provided)
	for i, suffix := range suffixes {
		if !strings.HasPrefix(suffix, ".") {
			suffixes[i] = "." + suffix
		}
	}

	var files []FileInfo

	// Walk through current directory and subdirectories
	err := filepath.Walk(".", func(path string, info os.FileInfo, err error) error {
		if err != nil {
			log.Printf("Error accessing %s: %v\n", path, err)
			return nil // Continue walking despite errors
		}

		// Skip directories
		if info.IsDir() {
			return nil
		}

		// Check if file has one of the target suffixes, or include all files if no suffixes specified
		if len(suffixes) == 0 {
			// No suffixes specified, include all files
			files = append(files, FileInfo{
				Path:    path,
				ModTime: info.ModTime(),
				Size:    info.Size(),
			})
		} else {
			// Check if file has one of the target suffixes
			for _, suffix := range suffixes {
				if strings.HasSuffix(strings.ToLower(info.Name()), strings.ToLower(suffix)) {
					files = append(files, FileInfo{
						Path:    path,
						ModTime: info.ModTime(),
						Size:    info.Size(),
					})
					break // Found a match, no need to check other suffixes
				}
			}
		}

		return nil
	})

	if err != nil {
		log.Printf("Error walking directory: %v\n", err)
		return
	}

	if len(files) == 0 {
		if len(suffixes) == 0 {
			fmt.Println("No files found.")
		} else {
			fmt.Println("No files found with the specified suffixes.")
		}
		return
	}

	// Sort files based on the selected option
	if *oldest {
		// Sort oldest to newest (ascending by modification time)
		sort.Slice(files, func(i, j int) bool {
			return files[i].ModTime.Before(files[j].ModTime)
		})
	} else if *largest {
		// Sort by largest files first (descending by size)
		sort.Slice(files, func(i, j int) bool {
			return files[i].Size > files[j].Size
		})
	} else if *smallest {
		// Sort by smallest files first (ascending by size)
		sort.Slice(files, func(i, j int) bool {
			return files[i].Size < files[j].Size
		})
	} else {
		// Default: Sort by newest first (descending by modification time)
		sort.Slice(files, func(i, j int) bool {
			return files[i].ModTime.After(files[j].ModTime)
		})
	}

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
