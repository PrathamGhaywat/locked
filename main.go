package main

import (
	"flag"
	"fmt"
	"os"
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
	fs.Parse(args)

	if fs.NArg() == 0 {
		fmt.Printf("Usage: locked lock <file or folder>")
		os.Exit(1)
	}

	path := fs.Arg(0)
	fmt.Printf("Locking: %s\n", path)
	fmt.Println("TODO")
}

func handleUnlock(args []string) {
	fs := flag.NewFlagSet("unlock", flag.ExitOnError)
	fs.Parse(args)

	if fs.NArg() == 0 {
		fmt.Println("Usage: locked unlock <file.locker>")
		os.Exit(1)
	}

	path := fs.Arg(0)
	fmt.Printf("Unlocking: %s\n", path)
	fmt.Println("TODO")
}

func printUsage() {
	fmt.Printf(`%s v%s - Ultra secure file locker
	
Usage:
%s lock <file or folder>	 Lock a file or folder
%s unlock <file.locker>		Unlock a .locker file
%s version							Show Version
%s help								  Show help message
`, AppName, Version, AppName, AppName,  AppName,  AppName)
}