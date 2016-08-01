package main

import (
	"fmt"
	"os"

	"gitlab.onegini.com/mobile-platform/onegini-sdk-configurator/cmd"
)

func main() {
	if err := cmd.RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
}
