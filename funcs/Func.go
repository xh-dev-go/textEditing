package funcs

import (
	"bufio"
	"fmt"
	"github.com/spf13/cobra"
	"io"
	"os"
)

type InputMode int

const (
	NewComplete InputMode = iota
	NewPartial
	RestingOfLine
)

type LineOut struct {
	InputMode
	String string
}

type StringProcessor interface {
	ProcessReader(
		reader *bufio.Reader,
		lineOut chan LineOut,
		done chan struct{},
	)
	Process(
		next func() ([]byte, bool, error),
		lineOut chan LineOut,
		done chan struct{},
	)
}

func NewStringProcessor(
	NewCompleteRow func(input string) string,
	NewPartialRow func(input string) string,
	RestingOfLine func(input string) string,
) StringProcessor {
	return &stringProcessorImpl{
		NewCompleteRow: NewCompleteRow,
		NewPartialRow:  NewPartialRow,
		RestingOfLine:  RestingOfLine,
	}
}

type stringProcessorImpl struct {
	NewCompleteRow func(input string) string
	NewPartialRow  func(input string) string
	RestingOfLine  func(input string) string
}

func (sp *stringProcessorImpl) ProcessReader(reader *bufio.Reader, lineOut chan LineOut, done chan struct{}) {
	sp.Process(
		func() ([]byte, bool, error) {
			return reader.ReadLine()
		},
		lineOut,
		done,
	)
}
func (sp *stringProcessorImpl) Process(next func() ([]byte, bool, error), lineOut chan LineOut, done chan struct{}) {
	var prefixOn = false
	for {
		lineBytes, isPrefix, err := next()
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
			l := sp.NewCompleteRow(string(lineBytes))
			lineOut <- LineOut{NewComplete, l}
		} else {

			if prefixOn {
				prefixOn = false
				l := sp.RestingOfLine(string(lineBytes))
				lineOut <- LineOut{RestingOfLine, l}
			} else {
				l := sp.NewPartialRow(string(lineBytes))
				lineOut <- LineOut{NewPartial, l}
			}
		}
	}
	done <- struct{}{}
}

func ForComplete[V any](msg chan V, done chan struct{}, operation func(input V)) {
	var isDone = false
	for {
		if isDone {
			break
		}
		select {
		case m := <-msg:
			operation(m)
		case <-done:
			isDone = true
		}
	}
}

func CommonExec(cmd *cobra.Command) {
	run, err := cmd.Flags().GetBool("exec")
	if err != nil {
		panic(err)
	}
	if !run {
		fmt.Println("No exec flag indicated!!!")
		fmt.Println()
		os.Exit(0)
	}
}
