package main

import (
	"fmt"
	"os"
)

var version = "dev"

func handleHelp(binName string) {
	options := Options{}
	handleVersion(binName)
	fmt.Println("A command line tool to handle unix timestamp")
	fmt.Println("")
	fmt.Println("USAGE:")
	fmt.Printf("  %s [OPTIONS] <SUBCOMMAND> [EXTRA_OPTIONS]\n", binName)
	fmt.Println("")

	fmt.Println("OPTIONS:")
	options.Flags().PrintOptions(os.Stdout)

	fmt.Println("")
	fmt.Println("SUBCOMMANDS:")
	fmt.Println("  generate   Generate unix timestamp with given options")
	fmt.Println("  help       Prints this message or the help of the given subcommand(s)")
	fmt.Println("  parse      Parse a unix timestamp and print it in human readable format")
}

func handleGenerateHelp(binName string) {
	options := GenerateOptions{}
	handleVersion(binName)
	fmt.Println("Generate unix timestamp with given options")
	fmt.Println("")
	fmt.Println("USAGE:")
	fmt.Printf("  %s [GENERAL_OPTIONS] generate [OPTIONS]\n", binName)
	fmt.Println("")
	fmt.Println("OPTIONS:")
	options.Flags().PrintOptions(os.Stdout)
}

func handleVersion(binName string) {
	fmt.Printf("%s %s\n", binName, version)
}

func run(runArgs ...string) error {
	var options Options

	binName := os.Args[0]
	args, err := options.Parse(runArgs...)
	if err != nil {
		return fmt.Errorf("%s: %w", runArgs[0], err)
	}

	if options.help != nil && *options.help {
		handleHelp(binName)
		return nil
	}

	if options.version != nil && *options.version {
		handleVersion(binName)
		return nil
	}

	if len(args) == 0 {
		return fmt.Errorf("%s: missing subcommand", runArgs[0])
	}

	switch args[0] {
	case "generate", "g":
		generateOptions := GenerateOptions{options: options}
		if remainingArgs, err := generateOptions.Parse(args...); err != nil {
			return fmt.Errorf("%s: %w", runArgs[0], err)
		} else if len(remainingArgs) > 0 {
			if remainingArgs[0] == "help" {
				handleGenerateHelp(binName)
				return nil
			}
			return fmt.Errorf("%s: unknown argument: %s", runArgs[0], remainingArgs[0])
		}
		if err := generate(os.Stdout, generateOptions); err != nil {
			return fmt.Errorf("%s: %w", runArgs[0], err)
		}
	case "parse", "p":
		if err := parse(os.Stdout, args[1:], options); err != nil {
			return fmt.Errorf("%s: %w", runArgs[0], err)
		}
	case "help", "h":
		handleHelp(binName)
		return nil
	default:
		return fmt.Errorf("%s: unknown subcommand: %s", runArgs[0], args[0])
	}

	return nil
}

func main() {
	if err := run(os.Args...); err != nil {
		fmt.Printf("%s: %s\n", os.Args[0], err)
		os.Exit(1)
	}
}
