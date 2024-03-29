/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"math"
	"os"
	"strings"

	"github.com/ross96D/mxargs/execute"
	"github.com/ross96D/mxargs/shared/config"
	"github.com/ross96D/mxargs/shared/flags"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

type flagType int

const (
	none flagType = iota
	short
	large
)

func parse(cmd *cobra.Command, args []string) ([]string, error) {
	i := 0
	for ; i < len(args); i++ {
		arg := args[i]
		var possibleNext string

		var flagtype flagType = none
		switch {
		case strings.HasPrefix(arg, "--"):
			arg, _ = strings.CutPrefix(arg, "--")

			splitted := strings.SplitN(arg, "=", 2)

			if len(splitted) == 2 {
				arg = splitted[0]
				possibleNext = splitted[1]
			}

			flagtype = large
		case strings.HasPrefix(arg, "-"):
			arg, _ = strings.CutPrefix(arg, "-")

			splitted := strings.SplitN(arg, "=", 2)

			if len(splitted) == 2 {
				arg = splitted[0]
				possibleNext = splitted[1]
			}

			flagtype = short
		}

		var flag *pflag.Flag
		switch flagtype {
		case short:
			flag = cmd.Flags().ShorthandLookup(arg)
		case large:
			flag = cmd.Flags().Lookup(arg)
		}

		if flag == nil {
			if flagtype != none {
				return nil, fmt.Errorf("%s is not a valid flag", arg)
			}
			break
		}

		switch flag.Value.Type() {
		case "int", "int32", "int64", "string", "uint", "uint64", "float64", "float32":
			if possibleNext != "" {
				arg = possibleNext
			} else {
				i += 1
				if i >= len(args) {
					return nil, fmt.Errorf("could not parse value for flag %s", flag.Name)
				}
				arg = args[i]
			}
			if err := flag.Value.Set(arg); err != nil {
				return nil, err
			}
		case "bool":
			flag.Value.Set("true")
		}
	}

	return args[i:], nil
}

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:                "mxargs [flags] [comand]",
	Short:              "A brief description of your application",
	Long:               ``,
	DisableFlagParsing: true,
	RunE: func(cmd *cobra.Command, args []string) (err error) {
		if args, err = parse(cmd, args); err != nil {
			return
		}
		// var mcmd *exec.Cmd
		// if len(args) > 1 {
		// 	mcmd = exec.Command(args[0], args[1:]...)
		// } else {
		// 	mcmd = exec.Command(args[0])
		// }
		// _ = mcmd
		conf, err := config.New(flags.Flags, os.Stdin)
		//! HANDLE ERROR
		if err != nil {
			return err
		}

		execute.ExecuteWithOrder(args, conf)

		return
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}
}

func init() {
	// rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

	flags := flags.Flags

	nullDesc := "Input items are terminated by the specified character." +
		"The specified delimiter may be a single character, a C-style character " +
		"escape such as \n, or an octal or hexadecimal escape code. " +
		"Octal and hexadecimal escape codes are understood as for the printf command. " +
		"Multibyte characters are not supported. When processing the input" +
		", quotes and backslash are not special; every character in the input is taken literally." +
		" The -d option disables any end-of-file string, which is treated like any other argument. " +
		"You can use this option when the input consists of simply newline-separated items, " +
		"although it is almost always better to design your program to use --null where this is possible."
	rootCmd.Flags().BoolVarP(&flags.NullSeparator, "null", "0", false, nullDesc)

	argFileDesc := "Read items from file instead of standard input. " +
		"If you use this option, stdin remains unchanged when commands are run. " +
		"Otherwise, stdin is redirected from /dev/null"
	rootCmd.Flags().StringVarP(&flags.FilePath, "arg-file", "a", "", argFileDesc)

	// TODO fix description
	delimeterDesc := "Input items are terminated by the specified character. " +
		"The specified delimiter may be a single character, a C-style character " +
		"escape such as \n, or an octal or hexadecimal escape code. " +
		"Octal and hexadecimal escape codes are understood as for the printf command. " +
		"Multibyte characters are not supported. When processing the input, quotes and " +
		"backslash are not special; every character in the input is taken literally. " +
		"The -d option disables any end-of-file string, which is treated like any other argument. " +
		"You can use this option when the input consists of simply newline-separated items, " +
		"although it is almost always better to design your program to use --null where " +
		"this is possible."
	rootCmd.Flags().StringVarP(&flags.Delimiter, "delimeter", "d", "\n", delimeterDesc)

	eofDesc := "Set the end of file string to eof-str. " +
		"If the end of file string occurs as a line of input, the rest of the input is ignored. " +
		"If -E is not used, no end of file string is used"
	rootCmd.Flags().StringVarP(&flags.EOF, "eof-str", "E", "", eofDesc)

	maxArgsDesc := "Use at most max-args arguments per command line. " +
		"Fewer than max-args arguments will be used if the size (see the -s option) is exceeded, " +
		"unless the -x option is given, in which case xargs will exit."
	rootCmd.Flags().IntVarP(&flags.MaxArgs, "max-args", "n", math.MaxInt, maxArgsDesc)

	sizeDesc := "Use at most max-chars characters per command line, including the command and " +
		"initial-arguments and the terminating nulls at the ends of the argument strings. " +
		"The largest allowed value is system-dependent, and is calculated as the argument " +
		"length limit for exec, less the size of your environment, less 2048 bytes of headroom. " +
		"If this value is more than 128KiB, 128KiB is used as the default value; otherwise, " +
		"the default value is the maximum. 1KiB is 1024 bytes."
	rootCmd.Flags().Uint64VarP(&flags.MaxChars, "max-chars", "s", 0, sizeDesc)

	verboseDesc := "Print the command line on the standard error output before executing it."
	rootCmd.Flags().BoolVarP(&flags.Verbose, "verbose", "t", false, verboseDesc)

	exitDesc := "Exit if the size (see the -s option) is exceeded."
	rootCmd.Flags().BoolVarP(&flags.Exit, "exit", "x", false, exitDesc)
}
