package cmd

import (
	"github.com/spf13/cobra"
	cloutpkg "github.com/tliron/puccini/clout"
	"github.com/tliron/puccini/clout/js"
	"github.com/tliron/puccini/common"
	formatpkg "github.com/tliron/puccini/common/format"
	"github.com/tliron/puccini/common/terminal"
	"github.com/tliron/puccini/tosca/compiler"
	urlpkg "github.com/tliron/puccini/url"
)

var output string
var resolve bool
var coerce bool
var exec string

func init() {
	rootCmd.AddCommand(compileCmd)
	compileCmd.Flags().StringArrayVarP(&inputs, "input", "i", []string{}, "specify an input (name=YAML)")
	compileCmd.Flags().StringVarP(&inputsUrl, "inputs", "n", "", "load inputs from a PATH or URL")
	compileCmd.Flags().StringVarP(&output, "output", "o", "", "output Clout to file (default is stdout)")
	compileCmd.Flags().BoolVarP(&resolve, "resolve", "r", true, "resolves the topology (attempts to satisfy all requirements with capabilities)")
	compileCmd.Flags().BoolVarP(&coerce, "coerce", "c", false, "coerces all values (calls functions and applies constraints)")
	compileCmd.Flags().StringVarP(&exec, "exec", "e", "", "execute JavaScript scriptlet")
}

var compileCmd = &cobra.Command{
	Use:   "compile [[TOSCA PATH or URL]]",
	Short: "Compile TOSCA to Clout",
	Long:  `Parses a TOSCA service template and compiles the normalized output of the parser to Clout. Supports JavaScript plugins.`,
	Args:  cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		var url string
		if len(args) == 1 {
			url = args[0]
		}

		Compile(url)
	},
}

func Compile(url string) {
	// Parse
	context, s := Parse(url)
	problems := context.GetProblems()

	// Compile
	clout, err := compiler.Compile(s)
	common.FailOnError(err)

	// Resolve
	if resolve {
		compiler.Resolve(clout, problems, format, pretty)
		FailOnProblems(problems)
	}

	// Coerce
	if coerce {
		compiler.Coerce(clout, problems, format, pretty)
		FailOnProblems(problems)
	}

	if exec != "" {
		err = Exec(exec, clout)
		common.FailOnError(err)
	} else if !terminal.Quiet || (output != "") {
		err = formatpkg.WriteOrPrint(clout, format, terminal.Stdout, pretty, output)
		common.FailOnError(err)
	}
}

func Exec(scriptletName string, clout *cloutpkg.Clout) error {
	clout, err := clout.Normalize()
	if err != nil {
		return err
	}

	// Try loading JavaScript from Clout
	scriptlet, err := js.GetScriptlet(scriptletName, clout)

	if err != nil {
		// Try loading JavaScript from path or URL
		url, err := urlpkg.NewValidURL(scriptletName, nil)
		common.FailOnError(err)

		scriptlet, err = urlpkg.Read(url)
		common.FailOnError(err)

		err = js.SetScriptlet(exec, js.CleanupScriptlet(scriptlet), clout)
		common.FailOnError(err)
	}

	jsContext := js.NewContext(scriptletName, log, terminal.Quiet, format, pretty, output)

	program, err := jsContext.GetProgram(scriptletName, scriptlet)
	if err != nil {
		return err
	}

	runtime := jsContext.NewCloutRuntime(clout, nil)

	_, err = runtime.RunProgram(program)

	return js.UnwrapException(err)
}
