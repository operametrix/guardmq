package cmd

import (
  "fmt"
  "github.com/spf13/cobra"
)

func init() {
  rootCmd.AddCommand(versionCmd)
}

var versionCmd = &cobra.Command{
  Use:   "version",
  Short: "Print the version number of GuardMQ",
  Long:  `All software has versions. This is GuardMQ's`,
  Run: func(cmd *cobra.Command, args []string) {
    fmt.Println("GuardMQ MQTT proxy for security and peering v1.0 -- HEAD")
  },
}