package cli

import (
	"errors"
	"fmt"
	"io"
	"strings"

	"github.com/chzyer/readline"

	"github.com/akansha204/cli-auth/internal/auth"
	"github.com/akansha204/cli-auth/internal/models"
	"github.com/akansha204/cli-auth/internal/session"
)

type App struct {
	auth     *auth.AuthService
	sessions *session.Manager
	rl       *readline.Instance

	user *models.User
	sess *models.Session
}

func NewApp(authService *auth.AuthService, sessionManager *session.Manager) *App {
	return &App{
		auth:     authService,
		sessions: sessionManager,
	}
}

func (a *App) Run() error {
	rl, err := readline.NewEx(&readline.Config{
		Prompt:       "> ",
		AutoComplete: newCompleter(),
	})
	if err != nil {
		return err
	}
	defer rl.Close()

	a.rl = rl

	for {
		line, err := rl.Readline()
		if err != nil {
			if errors.Is(err, readline.ErrInterrupt) {
				fmt.Println("(type 'quit' to exit)")
				continue
			}
			if errors.Is(err, io.EOF) {
				fmt.Println()
				return nil
			}
			return err
		}

		fields := strings.Fields(line)
		if len(fields) == 0 {
			continue
		}

		cmd, ok := lookup(fields[0])
		if !ok {
			fmt.Printf("unknown command: %s (type 'help')\n", fields[0])
			continue
		}

		if err := cmd.handler(a, fields[1:]); err != nil {
			if errors.Is(err, errQuit) {
				fmt.Println("bye")
				return nil
			}
			fmt.Printf("error: %s\n", err)
		}
	}
}

func (a *App) prompt(label string) (string, error) {
	a.rl.SetPrompt(label)

	line, err := a.rl.Readline()
	if err != nil {
		return "", err
	}

	return line, nil
}

func (a *App) promptPassword(label string) (string, error) {
	bytes, err := a.rl.ReadPassword(label)
	if err != nil {
		return "", err
	}

	return string(bytes), nil
}
