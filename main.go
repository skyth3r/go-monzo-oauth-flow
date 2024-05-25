package main

import (
	"fmt"
	"os"
)

const apiHostname = "api.monzo.com"

func main() {
	if err := run(); err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(1)
	}
}

func run() error {
	client := NewClient()

	if client.id == "" || client.secret == "" {
		return fmt.Errorf("the Client ID and Client secret were not found in env vars")
	}

	if client.accessToken == "" || client.refreshToken == "" {
		err := oauth(client)
		if err != nil {
			return err
		}
		os.Setenv("MONZO_ACCESS_TOKEN", client.accessToken)
		os.Setenv("MONZO_REFRESH_TOKEN", client.refreshToken)
	}

	err := pingTest(client)
	if err != nil {
		return err
	}

	fmt.Printf("Please open your Monzo app, click \"Allow access to your data\" for your application, and follow the instructions.\nOnce approved, press [Enter] to continue:\n")
	fmt.Scanln()

	// Example API call after SCA authentication
	err = accounts(client)
	if err != nil {
		return err
	}
	return nil
}
