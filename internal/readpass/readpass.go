package readpass

import (
	"errors"
	"os"

	"golang.org/x/term"
)

// PasswordPrompt is a simple (but echoing) password entry function
// that takes a prompt and reads the password.
func PasswordPrompt(prompt string) (password string, err error) {
	state, err := term.MakeRaw(int(os.Stdin.Fd()))
	if err != nil {
		return password, err
	}
	defer func() {
		errRestore := term.Restore(int(os.Stdin.Fd()), state)
		if errRestore != nil {
			if err == nil {
				err = errRestore
			} else {
				err = errors.Join(err, errRestore)
			}
		}
	}()

	t := term.NewTerminal(os.Stdout, ">")

	password, err = t.ReadPassword(prompt)
	if err != nil {
		return password, err
	}

	return password, err
}

// PasswordPromptBytes is a simple (but echoing) password entry function
// that takes a prompt and reads the password.
func PasswordPromptBytes(prompt string) (password []byte, err error) {
	passwordString, err := PasswordPrompt(prompt)
	if err != nil {
		return nil, err
	}

	password = []byte(passwordString)
	return
}
