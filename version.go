package main

import (
	"fmt"

	"github.com/spf13/cobra"
)

func versionCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "version",
		Short: "Print the version number of PortProxy",
		Long:  `All software has versions. This is PortProxy's`,
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("PortProxy v0.1.0")
		},
	}
}
