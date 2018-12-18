package main

import (
	"fmt"
	"os"
)

func main() {
	proxyCmd := proxyCmd()
	proxyCmd.AddCommand(versionCmd())

	if err := proxyCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
