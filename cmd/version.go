package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
)

const (
	version = "v0.0.1"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "打印当前软件版本信息",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("Version: %s\n", version)
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
