package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// versionCmd represents the version command
var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print version information",
	Long:  `Display the version of ResponseWatch CLI.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("ResponseWatch CLI (rwcli) version %s\n", version)
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
