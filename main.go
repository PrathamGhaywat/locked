package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/PrathamGhaywat/locked/pkg/vault"
)

const (
	Version = "0.1.0"
	AppName = "locked"
)

func main() {
	if len(os.Args) < 2 {
		printUsage()
		os.Exit(1)
	}

	command := os.Args[1]

	switch command {
	case "lock":
		handleLock(os.Args[2:])
	case "unlock":
		handleUnlock(os.Args[2:])
	case "version":
		fmt.Printf("%s v%s\n", AppName, Version)
	case "help":
		printUsage()
	default:
		fmt.Printf("Unknown command: %s\n", command)
		printUsage()
		os.Exit(1)
	}
}

func handleLock(args []string) {
	fs := flag.NewFlagSet("lock", flag.ExitOnError)
	fs.Usage = func() {
		fmt.Println("Usage: locked lock [options] <file or folder>")
		fmt.Println("\nOptions:")
		fs.PrintDefaults()
	}

	outputFlag := fs.String("o", "", "Output file (default: <input>.locker)")
	fs.Parse(args)

	if fs.NArg() == 0 {
		fs.Usage()
		os.Exit(1)
	}

	inputPath := fs.Arg(0)
	outputPath := *outputFlag

	// Set default output path if not specified
	if outputPath == "" {
		outputPath = inputPath + ".locker"
	}

	// Get password from user
	password, err := getPassword("Enter password: ")
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Locking: %s\n", inputPath)

	err = vault.CreateLocker(inputPath, outputPath, password)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("✓ Successfully locked to: %s\n", outputPath)
	fmt.Println("✓ Original file deleted")
}

func handleUnlock(args []string) {
	fs := flag.NewFlagSet("unlock", flag.ExitOnError)
	fs.Usage = func() {
		fmt.Println("Usage: locked unlock [options] <file.locker>")
		fmt.Println("\nOptions:")
		fs.PrintDefaults()
	}

	outputFlag := fs.String("o", "", "Output file (default: <original>_unlocked.<ext>)")
	fs.Parse(args)

	if fs.NArg() == 0 {
		fs.Usage()
		os.Exit(1)
	}

	lockerPath := fs.Arg(0)
	outputPath := *outputFlag

	// If no output path specified, extract original filename with "_unlocked" suffix
	if outputPath == "" {
		var err error
		originalFilename, err := vault.GetOriginalFilename(lockerPath)
		if err != nil {
			fmt.Printf("Error: %v\n", err)
			os.Exit(1)
		}

		// Add "_unlocked" before file extension
		ext := filepath.Ext(originalFilename)
		name := strings.TrimSuffix(originalFilename, ext)
		outputPath = name + "_unlocked" + ext
	}

	// Get password from user
	password, err := getPassword("Enter password: ")
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Unlocking: %s\n", lockerPath)

	err = vault.OpenLocker(lockerPath, outputPath, password)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("✓ Successfully unlocked to: %s\n", outputPath)
}

// getPassword reads a password from user input without echoing.
func getPassword(prompt string) (string, error) {
	fmt.Print(prompt)
	reader := bufio.NewReader(os.Stdin)
	password, err := reader.ReadString('\n')
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(password), nil
}

func printUsage() {
	fmt.Printf(`%s v%s - Ultra secure file locker

Usage:
  %s lock [options] <file>         Lock a file
  %s unlock [options] <file.locker> Unlock a .locker file
  %s version                        Show version
  %s help                           Show this help message

`, AppName, Version, AppName, AppName, AppName, AppName)
}
