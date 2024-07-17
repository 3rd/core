package cmd

import (
	"core/utils"
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var cfgFile string
var env = utils.GetEnv()

var rootCmd = &cobra.Command{
	Use:   "core",
	Short: "core",
	Long:  "core",
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize()
}
