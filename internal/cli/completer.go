package cli

import "github.com/chzyer/readline"

func newCompleter() *readline.PrefixCompleter {
	items := make([]readline.PrefixCompleterInterface, 0, len(commandList()))
	for _, c := range commandList() {
		items = append(items, readline.PcItem(c.name))
	}

	return readline.NewPrefixCompleter(items...)
}
