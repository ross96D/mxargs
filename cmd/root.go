/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"sync"

	"github.com/ross96D/mxargs/execute"
	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "mxargs",
	Short: "A brief description of your application",
	Long: `A longer description that spans multiple lines and likely contains
examples and usage of using your application. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) == 0 {
			return fmt.Errorf("must be at least 1 argument")
		}
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		var mcmd *exec.Cmd
		if len(args) > 1 {
			mcmd = exec.Command(args[0], args[1:]...)
		} else {
			mcmd = exec.Command(args[0])
		}
		_ = mcmd

		var p = execute.Print{
			Mut: &sync.Mutex{},
		}
		argsss := stdin()
		execute.Execute(&p, mcmd, argsss)
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	// rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.mxargs.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func stdin() []string {
	stdin := os.Stdin
	b, err := io.ReadAll(stdin)
	if err != nil {
		return []string{}
	}
	resp := make([]string, 0)

	buff := make([]byte, 0, len(b))
	for i := 0; i < len(b); i++ {
		if b[i] == 0 {
			resp = append(resp, string(buff))
			buff = buff[0:0]
			continue
		}
		buff = append(buff, b[i])
	}
	return resp
}
