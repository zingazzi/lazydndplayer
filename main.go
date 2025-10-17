// main.go
package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/marcozingoni/lazydndplayer/internal/storage"
	"github.com/marcozingoni/lazydndplayer/internal/ui"
)

func main() {
	// Parse command line flags
	charFile := flag.String("file", storage.GetDefaultPath(), "Path to character file")
	importFile := flag.String("import", "", "Import character from file")
	exportFile := flag.String("export", "", "Export character to file")
	flag.Parse()

	// Initialize storage
	store := storage.NewStorage(*charFile)

	// Handle import
	if *importFile != "" {
		char, err := store.Import(*importFile)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error importing character: %v\n", err)
			os.Exit(1)
		}

		err = store.Save(char)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error saving character: %v\n", err)
			os.Exit(1)
		}

		fmt.Printf("Character imported successfully to %s\n", *charFile)
		return
	}

	// Load character
	char, err := store.Load()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error loading character: %v\n", err)
		os.Exit(1)
	}

	// Handle export
	if *exportFile != "" {
		err := store.Export(char, *exportFile)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error exporting character: %v\n", err)
			os.Exit(1)
		}
		fmt.Printf("Character exported successfully to %s\n", *exportFile)
		return
	}

	// Run the TUI
	if err := ui.Run(char, store); err != nil {
		fmt.Fprintf(os.Stderr, "Error running application: %v\n", err)
		os.Exit(1)
	}
}
