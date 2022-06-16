/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"github.com/spf13/cobra"
	"github.com/xh-dev-go/textEditing/cmd/line"
)

// ineCmd represents the line command
var lineCmd = &cobra.Command{
	Use:   "line",
	Short: "Editing text line by line",
	Long:  `Editing text line by line`,
}

func init() {
	rootCmd.AddCommand(lineCmd)
	lineCmd.AddCommand(line.LnCmd)
	lineCmd.AddCommand(line.RmCmd)
	lineCmd.PersistentFlags().BoolP("exec", "e", false, "execute the command")

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// ineCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// ineCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
