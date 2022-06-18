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
	"io"
	"os"
)

// rmCmd represents the rm command
var RmCmd = &cobra.Command{
	Use:   "tail",
	Short: "keep the tail of string",
	Long: `Split each row by column index or separator and keep the tail

Input: 
sep->xx
msg->sdfxxsdfasdxxx\nddddddxxjsdkfjs

Output:
sdfasdxxx\njsdkfjs

`,
	Run: func(cmd *cobra.Command, args []string) {
		funcs.CommonExec(cmd)

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
		sep, err := cmd.Flags().GetString("separator")
		if err != nil {
			panic(err)
		}
		col, err := cmd.Flags().GetInt("column")
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
				if prefixOn {
					l := string(lineBytes)
					_, err := fmt.Fprintf(os.Stdout, l)
					if err != nil {
						panic(err)
					}
					if cbOut {
						msgOut += l
					}
				} else {
					prefixOn = true
					index++
					var l string
					if col == 0 {
						var replaced bool
						l, replaced = stringUtils.RemoveFirst(string(lineBytes), sep)
						if !replaced {
							panic("No replace for: " + string(lineBytes))
						}
					} else {
						l = string(lineBytes)[col:]
					}

					_, err := fmt.Fprintf(os.Stdout, l)
					if err != nil {
						panic(err)
					}
					if cbOut {
						msgOut += l
					}
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
					var l string
					if col == 0 {
						var replaced bool
						l, replaced = stringUtils.RemoveFirst(string(lineBytes), sep)
						if !replaced {
							panic("No replace for: " + string(lineBytes))
						}
					} else {
						l = string(lineBytes)[col:]
					}
					l = l + "\n"
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
			panic(err)
		}

	},
}

func init() {
	RmCmd.PersistentFlags().IntP("column", "c", 0, "keep tail from column index")
	RmCmd.PersistentFlags().StringP("separator", "s", "\t", "separator to be removed.")
	RmCmd.PersistentFlags().BoolP("clipboard-in", "i", false, "The source is from clipboard, otherwise from the stdin")
	RmCmd.PersistentFlags().BoolP("clipboard-out", "o", false, "The destination to clipboard also, default to stdout only")

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// rmCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// rmCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
