package cli

import "errors"

var errExit = errors.New("exit")

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
		{name: "whoami", usage: "show the current session", handler: cmdStatus},
		{name: "enable-2fa", usage: "enable two-factor authentication", handler: cmdEnableMFA},
		{name: "disable-2fa", usage: "disable two-factor authentication", handler: cmdDisableMFA},
		{name: "exit", usage: "exit the application", handler: cmdExit},
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
