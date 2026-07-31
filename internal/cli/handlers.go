package cli

import (
	"errors"
	"fmt"
	"strings"

	"github.com/akansha204/cli-auth/internal/auth"
	"github.com/akansha204/cli-auth/internal/session"
	"github.com/akansha204/cli-auth/internal/utils"
)

func cmdHelp(a *App, _ []string) error {
	for _, c := range commandList() {
		fmt.Printf("  %-12s %s\n", c.name, c.usage)
	}
	return nil
}

func cmdRegister(a *App, _ []string) error {
	username, err := a.prompt("Username: ")
	if err != nil {
		return err
	}

	password, err := a.promptPassword("Password: ")
	if err != nil {
		return err
	}

	confirm, err := a.promptPassword("Confirm password: ")
	if err != nil {
		return err
	}

	if strings.TrimSpace(password) != strings.TrimSpace(confirm) {
		return errors.New("passwords do not match")
	}

	if err := a.auth.Register(username, password); err != nil {
		return err
	}

	fmt.Printf("Registered %s. You can now log in.\n", strings.ToLower(strings.TrimSpace(username)))
	return nil
}

func cmdLogin(a *App, _ []string) error {
	if a.user != nil {
		return fmt.Errorf("already logged in as %s", a.user.Username)
	}

	username, err := a.prompt("Username: ")
	if err != nil {
		return err
	}

	password, err := a.promptPassword("Password: ")
	if err != nil {
		return err
	}

	user, err := a.auth.Login(username, password)
	if errors.Is(err, auth.ErrMFARequired) {
		code, promptErr := a.prompt("MFA code: ")
		if promptErr != nil {
			return promptErr
		}

		user, err = a.auth.VerifyMFA(username, code)
	}
	if err != nil {
		return err
	}

	sess, err := a.sessions.Create(user.ID)
	if err != nil {
		return err
	}

	a.user = user
	a.sess = sess

	return printSessionInfo(a)
}

func cmdLogout(a *App, _ []string) error {
	if a.user == nil {
		return errors.New("not logged in")
	}

	if err := a.sessions.Invalidate(a.user.ID); err != nil && !errors.Is(err, session.ErrNoSession) {
		return err
	}

	a.user = nil
	a.sess = nil

	fmt.Println("Logged out")
	return nil
}

func cmdStatus(a *App, _ []string) error {
	if a.user == nil {
		fmt.Println("Not logged in")
		return nil
	}

	return printSessionInfo(a)
}

func printSessionInfo(a *App) error {
	sess, err := a.sessions.Active(a.user.ID)
	if errors.Is(err, session.ErrNoSession) {
		a.user = nil
		a.sess = nil
		return errors.New("session is no longer active, please log in again")
	}
	if errors.Is(err, session.ErrSessionExpired) {
		a.user = nil
		a.sess = nil
		return errors.New("session has expired, please log in again")
	}
	if err != nil {
		return err
	}

	a.sess = sess

	fmt.Printf("Logged in as %s\n", a.user.Username)
	fmt.Printf("Registered: %s\n", utils.FormatTime(a.user.RegisteredAt))
	fmt.Printf("Session expires: %s\n", utils.FormatTime(sess.ExpiresAt))

	if a.user.LastLogin != nil {
		fmt.Printf("Last login: %s\n", utils.FormatTime(*a.user.LastLogin))
	}

	mfa := "disabled"
	if a.user.MFAEnabled {
		mfa = "enabled"
	}
	fmt.Printf("MFA: %s\n", mfa)

	return nil
}

func cmdEnableMFA(a *App, _ []string) error {
	if a.user == nil {
		return errors.New("not logged in")
	}

	secret, uri, err := a.auth.EnableMFA(a.user.Username)
	if err != nil {
		return err
	}

	a.user.MFAEnabled = true

	fmt.Println("MFA enabled. Add this account to your authenticator app:")
	fmt.Printf("  Secret: %s\n", secret)
	fmt.Printf("  URI:    %s\n", uri)
	return nil
}

func cmdDisableMFA(a *App, _ []string) error {
	if a.user == nil {
		return errors.New("not logged in")
	}

	if err := a.auth.DisableMFA(a.user.Username); err != nil {
		return err
	}

	a.user.MFAEnabled = false

	fmt.Println("MFA disabled")
	return nil
}

func cmdExit(_ *App, _ []string) error {
	return errExit
}
