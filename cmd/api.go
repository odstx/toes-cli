package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
)

var apiCmd = &cobra.Command{
	Use:   "api",
	Short: "全新完全创建一个 api demo 项目",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("完整生成 api demo 项目")
	},
}

func init() {
	rootCmd.AddCommand(apiCmd)
}
