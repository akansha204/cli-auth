package cli

import "errors"

var errQuit = errors.New("quit")

type command struct {
	name    string
	usage   string
	handler func(*App, []string) error
}

func commandList() []command {
	return []command{
		{name: "help", usage: "show available commands", handler: cmdHelp},
		{name: "register", usage: "create a new account", handler: cmdRegister},
		{name: "login", usage: "log in to your account", handler: cmdLogin},
		{name: "logout", usage: "end the current session", handler: cmdLogout},
		{name: "status", usage: "show the current session", handler: cmdStatus},
		{name: "mfa", usage: "enable two-factor authentication", handler: cmdEnableMFA},
		{name: "disable-mfa", usage: "disable two-factor authentication", handler: cmdDisableMFA},
		{name: "quit", usage: "exit the application", handler: cmdQuit},
	}
}

func lookup(name string) (*command, bool) {
	for _, c := range commandList() {
		if c.name == name {
			return &c, true
		}
	}
	return nil, false
}
