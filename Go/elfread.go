package main

import (
	"debug/elf"
	"fmt"
	"log"
	"os"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run main.go <elf-file>")
		os.Exit(1)
	}

	filePath := os.Args[1]
	file, err := elf.Open(filePath)
	if err != nil {
		log.Fatalf("Failed to open ELF file: %v", err)
	}
	defer file.Close()

	// Get the dynamic symbols which include imported functions
	symbols, err := file.DynamicSymbols()
	if err != nil {
		log.Fatalf("Failed to read dynamic symbols: %v", err)
	}

	importedLibraries := make(map[string]map[string]struct{})

	for _, symbol := range symbols {
		if symbol.Library != "" {
			if _, ok := importedLibraries[symbol.Library]; !ok {
				importedLibraries[symbol.Library] = make(map[string]struct{})
			}
			importedLibraries[symbol.Library][symbol.Name] = struct{}{}
		}
	}

	fmt.Println("Imports and their functions:")
	for lib, funcs := range importedLibraries {
		fmt.Printf("Library: %s\n", lib)
		for funcName := range funcs {
			fmt.Printf("  Function: %s\n", funcName)
		}
	}
}
