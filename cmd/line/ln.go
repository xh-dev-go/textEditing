/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package line

import (
	"bufio"
	"fmt"
	"github.com/atotto/clipboard"
	"github.com/spf13/cobra"
	"github.com/xh-dev-go/textEditing/funcs"
	"github.com/xh-dev-go/xhUtils/stringUtils"
	"os"
)

// LnCmd represents the add command
var LnCmd = &cobra.Command{
	Use:     "ln",
	Aliases: []string{"line-num"},
	Short:   "Add a line number to every row",
	Long:    `Add a line number to every row`,
	Run: func(cmd *cobra.Command, args []string) {
		funcs.CommonExec(cmd)

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

		var index = 0

		var out = make(chan funcs.LineOut)
		var done = make(chan struct{})
		func() {
			processor := funcs.NewStringProcessor(
				func(input string) string {
					index++
					return stringUtils.StringConcat(stringUtils.FormatNum(index, format), sep, input, false)
				},
				func(input string) string {
					index++
					return stringUtils.StringConcat(stringUtils.FormatNum(index, format), sep, input, true)
				},
				func(input string) string {
					return fmt.Sprintf("%s", input)
				},
			)
			go processor.ProcessReader(
				reader,
				out,
				done,
			)
		}()

		var msgOut = ""
		funcs.ForComplete(out, done, func(lineOut funcs.LineOut) {
			switch lineOut.InputMode {
			case funcs.NewComplete:
				fmt.Print(lineOut.String)
				if cbOut {
					msgOut += lineOut.String
				}
			case funcs.NewPartial:
				fmt.Print(lineOut.String)
				if cbOut {
					msgOut += lineOut.String
				}
			case funcs.RestingOfLine:
				fmt.Print(lineOut.String)
				if cbOut {
					msgOut += lineOut.String
				}
			}
		})

		if cbOut {
			err = clipboard.WriteAll(msgOut)
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
