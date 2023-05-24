package main

import (
	"fmt"

	"github.com/bradmccoydev/cdevents-controller/pkg/version"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(versionCmd)
}

var versionCmd = &cobra.Command{
	Use:   `version`,
	Short: "Prints cdeventscli version",
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Println(version.VERSION)
		return nil
	},
}
