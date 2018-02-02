package main

import (
	"fmt"
	"os"
	"strings"

	prompt "github.com/c-bata/go-prompt"
	colorable "github.com/mattn/go-colorable"
)

const (
	helpTemplate = "This is a wrapper around the 'snowy' repository."
)

var list = commands{
	suggestion("help", "print help"),
	suggestion("health", "show status",
		flag("--status", "status only health check"),
	),
	suggestion("ledger", "show a ledger"),
	suggestion("ledgers", "show ledgers as revisions"),
}

func main() {
	p := prompt.New(
		execute(client{
			base: "http://0.0.0.0:8080",
			out:  colorable.NewColorableStdout(),
		}),
		complete(list),
		prompt.OptionTitle("snowy: interactive Snowy client (type \"q\" to exit)"),
		prompt.OptionPrefix("> "),
		prompt.OptionInputTextColor(prompt.Yellow),
	)
	fmt.Fprintln(os.Stdout, "snowy (type \"q\" to quit)")
	p.Run()
}

func execute(c client) func(string) {
	return func(s string) {
		switch a := strings.TrimSpace(s); a {
		case "":
			return
		case "?", "h":
			fmt.Fprintln(os.Stdout, helpTemplate)
		case "quit", "exit", "q":
			fmt.Fprintln(os.Stdout, "Bye!")
			os.Exit(0)
		default:
			execRequest(c, a)
		}
	}
}

func execRequest(c client, s string) {
	parts := strings.Split(s, " ")
	switch parts[0] {
	case "help":
		fmt.Fprintln(os.Stdout, helpTemplate)
		os.Exit(0)
	case "health":
		statusOnly := false
		validateOptions(parts[1:], func(s string) bool {
			switch s {
			case "--status":
				statusOnly = true
				return true
			default:
				return false
			}
		})
		c.health(statusOnly)
	case "ledger":
		validateOptions(parts[1:], func(s string) bool {
			return !strings.HasPrefix(s, "-")
		})
		c.ledger(parts[1])
	case "ledgers":
		validateOptions(parts[1:], func(s string) bool {
			return !strings.HasPrefix(s, "-")
		})
		c.ledgers(parts[1])
	}
}

func complete(commands commands) func(prompt.Document) []prompt.Suggest {
	return func(d prompt.Document) []prompt.Suggest {
		var suggestions []prompt.Suggest

		text := d.TextBeforeCursor()
		if text == "" || strings.Contains(text, "|") {
			return suggestions
		}

		args := strings.Split(text, " ")

		return completeSubcommands(commands, args)
	}
}

func completeSubcommands(commands commands, args []string) []prompt.Suggest {
	if len(args) < 2 {
		return prompt.FilterHasPrefix(commands.suggestions(), args[0], true)
	}

	subCommands, ok := commands.get(args[0])
	if !ok {
		return []prompt.Suggest{}
	}

	return subCommands.suggestFor(args[1:])
}

func validateOptions(options []string, fn func(string) bool) {
	for _, v := range options {
		if !fn(v) {
			os.Stderr.WriteString(fmt.Sprintf("Invalid argument %q\n", v))
			os.Exit(1)
		}
	}
}
