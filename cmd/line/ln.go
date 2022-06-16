/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package line

import (
	"bufio"
	"fmt"
	"github.com/atotto/clipboard"
	"github.com/spf13/cobra"
	"github.com/xh-dev-go/xhUtils/stringUtils"
	"io"
	"os"
)

// LnCmd represents the add command
var LnCmd = &cobra.Command{
	Use:     "ln",
	Aliases: []string{"line-num"},
	Short:   "Add a line number to every row",
	Long:    `Add a line number to every row`,
	Run: func(cmd *cobra.Command, args []string) {
		run, err := cmd.Flags().GetBool("exec")
		if err != nil {
			if !run {
				fmt.Println("No exec flag indicated!!!")
				fmt.Println()
			}
			cmd.Usage()
			return
		}

		var index = 0
		var prefixOn = false

		cbIn, err := cmd.LocalFlags().GetBool("clipboard-in")
		if err != nil {
			panic(err)
		}
		cbOut, err := cmd.LocalFlags().GetBool("clipboard-out")
		if err != nil {
			panic(err)
		}
		sep, err := cmd.LocalFlags().GetString("separator")
		if err != nil {
			panic(err)
		}
		format, err := cmd.LocalFlags().GetString("index-format")
		if err != nil {
			panic(err)
		}

		var reader *bufio.Reader
		if cbIn {
			strs, err := clipboard.ReadAll()
			if err != nil {
				panic(err)
			}
			reader = stringUtils.StringToBufIoReader(strs)
		} else {
			reader = bufio.NewReader(os.Stdin)
		}

		var msgOut = ""
		for {
			lineBytes, isPrefix, err := reader.ReadLine()
			if err != nil {

				// Complete task
				if err == io.EOF {
					break
				} else {
					_, err := fmt.Fprintf(os.Stderr, err.Error())
					if err != nil {
						panic(err)
					}
				}
			}

			if isPrefix {
				prefixOn = true
				index++
				l := stringUtils.StringConcat(stringUtils.FormatNum(index, format), sep, string(lineBytes), false)

				_, err := fmt.Fprintf(os.Stdout, l)
				if err != nil {
					panic(err)
				}
				if cbOut {
					msgOut += l
				}
			} else {

				if prefixOn {
					prefixOn = false
					l := fmt.Sprintf("%s\n", string(lineBytes))
					_, err := fmt.Fprintf(os.Stdout, l)
					if err != nil {
						panic(err)
					}
					if cbOut {
						msgOut += l
					}
				} else {
					index++
					l := stringUtils.StringConcat(stringUtils.FormatNum(index, format), sep, string(lineBytes), true)
					_, err := fmt.Fprintf(os.Stdout, l)
					if err != nil {
						panic(err)
					}
					if cbOut {
						msgOut += l
					}
				}
			}
		}
		if cbOut {
			err = clipboard.WriteAll(msgOut)
		}
		if err != nil {
			return
		}
	},
}

func init() {

	LnCmd.PersistentFlags().StringP("index-format", "f", "%d", "The formatting for line number.")
	LnCmd.PersistentFlags().StringP("separator", "s", "\t", "The separator between line number and row. [Default '\\n']")
	LnCmd.PersistentFlags().BoolP("clipboard-in", "i", false, "The source is from clipboard, otherwise from the stdin")
	LnCmd.PersistentFlags().BoolP("clipboard-out", "o", false, "The destination to clipboard also, default to stdout only")

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// LnCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// LnCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
