/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"github.com/spf13/cobra"
	"github.com/starjun/toes-cli/internal/dev"
)

// devCmd represents the dev command
var devCmd = &cobra.Command{
	Use:   "dev",
	Short: "启动一个dev环境 redis+mysql",
	Long:  `快速启动一个开发环境`,
	Run: func(cmd *cobra.Command, args []string) {
		dev.Run(args)
	},
}

func init() {
	rootCmd.AddCommand(devCmd)
}
