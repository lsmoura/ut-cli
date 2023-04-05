package main

import (
	"fmt"
	"os"
)

var version = "dev"

func handleHelp() {
	options := Options{}
	handleVersion()
	fmt.Println("A command line tool to handle unix timestamp")
	fmt.Println("")
	fmt.Println("USAGE:")
	fmt.Println("  ut [OPTIONS] <SUBCOMMAND> [EXTRA_OPTIONS]")
	fmt.Println("")

	fmt.Println("OPTIONS:")
	options.Flags().PrintOptions(os.Stdout)

	fmt.Println("")
	fmt.Println("SUBCOMMANDS:")
	fmt.Println("  generate   Generate unix timestamp with given options")
	fmt.Println("  help       Prints this message or the help of the given subcommand(s)")
	fmt.Println("  parse      Parse a unix timestamp and print it in human readable format")
}

func handleGenerateHelp() {
	options := GenerateOptions{}
	handleVersion()
	fmt.Println("Generate unix timestamp with given options")
	fmt.Println("")
	fmt.Println("USAGE:")
	fmt.Println("  ut [GENERAL_OPTIONS] generate [OPTIONS]")
	fmt.Println("")
	fmt.Println("OPTIONS:")
	options.Flags().PrintOptions(os.Stdout)
}

func handleVersion() {
	fmt.Printf("ut %s\n", version)
}

func main() {
	var options Options

	args, err := options.Parse(os.Args...)
	if err != nil {
		fmt.Printf("%s: %s\n", os.Args[0], err)
		os.Exit(1)
	}

	if options.help != nil && *options.help {
		handleHelp()
		os.Exit(0)
	}

	if options.version != nil && *options.version {
		handleVersion()
		os.Exit(0)
	}

	if len(args) == 0 {
		fmt.Printf("%s: missing subcommand\n", os.Args[0])
		os.Exit(1)
	}

	switch args[0] {
	case "generate", "g":
		generateOptions := GenerateOptions{options: options}
		if remainingArgs, err := generateOptions.Parse(args...); err != nil {
			fmt.Printf("%s: %s\n", os.Args[0], err)
			os.Exit(1)
		} else if len(remainingArgs) > 0 {
			if remainingArgs[0] == "help" {
				handleGenerateHelp()
				os.Exit(0)
			}
			fmt.Printf("%s: unknown argument: %s\n", os.Args[0], remainingArgs[0])
			os.Exit(1)
		}
		if err := generate(os.Stdout, generateOptions); err != nil {
			fmt.Printf("%s: %s\n", os.Args[0], err)
			os.Exit(1)
		}
	case "parse", "p":
		if err := parse(os.Stdout, args[1:], options); err != nil {
			fmt.Printf("%s: %s\n", os.Args[0], err)
			os.Exit(1)
		}
	case "help", "h":
		handleHelp()
	default:
		fmt.Printf("%s: unknown subcommand: %s\n", os.Args[0], args[0])
		os.Exit(1)
	}
}
