package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
)

// convertWorkshops converts all workshop HTML files to PNG
func convertWorkshops() {
	// Dynamically find all workshop HTML files
	workshopDir := "../open-workshops"
	outputDir := "../open-workshops/output"
	pattern := filepath.Join(workshopDir, "workshop-*.html")

	// Convert output directory to absolute path
	absOutputDir, err := filepath.Abs(outputDir)
	if err != nil {
		log.Fatalf("Failed to get absolute path for output directory: %v", err)
	}

	// Create output directory for PNG files
	if err := os.MkdirAll(absOutputDir, 0755); err != nil {
		log.Fatalf("Failed to create output directory: %v", err)
	}

	matches, err := filepath.Glob(pattern)
	if err != nil {
		log.Fatalf("Failed to search for workshop files: %v", err)
	}

	if len(matches) == 0 {
		log.Printf("No workshop HTML files found in %s", workshopDir)
		log.Printf("Make sure to run the HTML generation first: cd ../open-workshops && go run main.go")
		return
	}

	fmt.Printf("Found %d workshop HTML files to convert...\n", len(matches))
	fmt.Printf("PNG files will be saved to: %s\n", absOutputDir)

	for _, htmlFile := range matches {
		// Check if HTML file exists (double-check)
		if _, err := os.Stat(htmlFile); os.IsNotExist(err) {
			log.Printf("HTML file not found: %s", htmlFile)
			continue
		}

		fmt.Printf("Converting %s...\n", filepath.Base(htmlFile))

		// Convert relative path to absolute path
		absPath, err := filepath.Abs(htmlFile)
		if err != nil {
			log.Printf("Failed to get absolute path for %s: %v", htmlFile, err)
			continue
		}

		// Convert to file:// URL format
		var fileURL string
		if runtime.GOOS == "windows" {
			// Windows: file:///C:/path/to/file.html
			fileURL = "file:///" + absPath
		} else {
			// Unix-like: file:///path/to/file.html
			fileURL = "file://" + absPath
		}

		// Change working directory to html-to-png for the conversion
		// but specify output directory for PNG files
		cmd := exec.Command("go", "run", "main.go", fileURL)
		cmd.Dir = "."
		cmd.Env = append(os.Environ(), "OUTPUT_DIR="+absOutputDir)

		output, err := cmd.CombinedOutput()

		if err != nil {
			log.Printf("Failed to convert %s to PNG: %v\nOutput: %s", htmlFile, err, string(output))
		} else {
			fmt.Printf("✅ Successfully converted: %s\n", filepath.Base(htmlFile))
		}
	}

	fmt.Println("\nAll workshop conversions completed!")
	fmt.Printf("PNG files saved to: %s\n", absOutputDir)
}

func main() {
	convertWorkshops()
}
