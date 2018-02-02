package main

import prompt "github.com/c-bata/go-prompt"
import "strings"

const (
	defaultOptionPrefix = "-"
)

type command struct {
	suggest  prompt.Suggest
	commands commands
	options  options
}

func (c command) suggestFor(args []string) (res []prompt.Suggest) {
	if len(args) == 0 {
		return
	}

	if a := args[len(args)-1]; strings.HasPrefix(a, defaultOptionPrefix) {
		return prompt.FilterHasPrefix(c.options.suggestions(), a, true)
	}

	suggestions := c.commands.suggestions()
	return prompt.FilterHasPrefix(suggestions, args[0], true)
}

type commands []command

func (c commands) suggestions() []prompt.Suggest {
	res := make([]prompt.Suggest, len(c))
	for k, v := range c {
		res[k] = v.suggest
	}
	return res
}

func (c commands) get(name string) (command, bool) {
	for _, v := range c {
		if v.suggest.Text == name {
			return v, true
		}
	}
	return command{}, false
}

type option struct {
	suggest prompt.Suggest
}

type options []option

func (o options) suggestions() []prompt.Suggest {
	res := make([]prompt.Suggest, len(o))
	for k, v := range o {
		res[k] = v.suggest
	}
	return res
}

func suggestion(text, desc string, opts ...option) command {
	return command{
		suggest: prompt.Suggest{
			Text:        text,
			Description: desc,
		},
		options: options(opts),
	}
}

func flag(text, desc string) option {
	return option{
		suggest: prompt.Suggest{
			Text:        text,
			Description: desc,
		},
	}
}
