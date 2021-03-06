package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/tliron/puccini/common"
	"github.com/tliron/puccini/common/terminal"
)

var logTo string
var verbose int
var format string
var colorize bool
var pretty bool
var problemsFormat string
var quirks []string

var bashCompletionTo string

func init() {
	rootCmd.PersistentFlags().BoolVarP(&terminal.Quiet, "quiet", "q", false, "suppress output")
	rootCmd.PersistentFlags().StringVarP(&logTo, "log", "l", "", "log to file (defaults to stderr)")
	rootCmd.PersistentFlags().CountVarP(&verbose, "verbose", "v", "add a log verbosity level (can be used twice)")
	rootCmd.PersistentFlags().StringVarP(&format, "format", "f", "", "force output format (\"yaml\", \"json\", or \"xml\")")
	rootCmd.PersistentFlags().BoolVarP(&colorize, "colorize", "z", true, "colorize output")
	rootCmd.PersistentFlags().BoolVarP(&pretty, "pretty", "p", true, "prettify output")
	rootCmd.PersistentFlags().StringVarP(&problemsFormat, "problems-format", "m", "", "problems format (\"yaml\", \"json\", or \"xml\")")
	rootCmd.PersistentFlags().StringSliceVarP(&quirks, "quirk", "x", nil, "parser quirk")

	rootCmd.Flags().StringVarP(&bashCompletionTo, "bash-completion", "b", "", "generate bash completion file")
}

var rootCmd = &cobra.Command{
	Use:   "puccini-tosca",
	Short: "TOSCA frontend for Clout",
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		if colorize {
			terminal.EnableColor()
		}
		if logTo == "" {
			if terminal.Quiet {
				verbose = -4
			}
			common.ConfigureLogging(verbose, nil)
		} else {
			common.ConfigureLogging(verbose, &logTo)
		}
	},
	Run: func(cmd *cobra.Command, args []string) {
		if bashCompletionTo != "" {
			if !terminal.Quiet {
				fmt.Fprintf(terminal.Stdout, "generating bash completion script: %s\n", bashCompletionTo)
			}
			cmd.GenBashCompletionFile(bashCompletionTo)
		}
	},
}

func Execute() {
	err := rootCmd.Execute()
	common.FailOnError(err)
}
