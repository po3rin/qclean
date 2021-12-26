package main

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/po3rin/qclean"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"golang.org/x/crypto/ssh/terminal"
)

var rootCmd = &cobra.Command{
	Use:   "qclean",
	Short: "qclean lets you to clean up search query in japanese. This is mainly used to remove wasted space.",
	Run: func(cmd *cobra.Command, args []string) {
		var c *qclean.Cleaner
		var err error

		if terminal.IsTerminal(int(os.Stdin.Fd())) {
			fmt.Fprintln(os.Stdout, "Currently does not support interactive mode")
			os.Exit(1)
		}

		c, err = qclean.NewCleaner()
		if err != nil {
			fmt.Fprintln(os.Stdout, err)
			os.Exit(1)
		}

		input, err := ioutil.ReadAll(os.Stdin)
		if err != nil {
			fmt.Fprintln(os.Stdout, err.Error())
			os.Exit(1)
		}

		result, err := c.Clean(string(input))
		if err != nil {
			fmt.Fprintln(os.Stdout, err)
			os.Exit(1)
		}
		fmt.Println(result)
	},
}

func init() {
	viper.SetEnvPrefix("qclean")
	viper.AutomaticEnv()

	viper.BindPFlags(pflag.CommandLine)
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stdout, err)
		os.Exit(1)
	}
}
