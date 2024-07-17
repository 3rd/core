package cmd

import (
	"fmt"
	"time"

	wiki "github.com/3rd/core/core-lib/wiki/local"
	"github.com/spf13/cobra"
)

var currentTaskCommand = &cobra.Command{
	Use:   "current",
	Short: "list the currently in-progress task (first only)",
	Run: func(cmd *cobra.Command, args []string) {
		root := env.WIKI_ROOT
		if len(root) == 0 {
			panic("WIKI_ROOT not set")
		}

		wiki, err := wiki.NewLocalWiki(wiki.LocalWikiConfig{
			Root:  root,
			Parse: "full",
		})
		if err != nil {
			panic(err)
		}

		nodes, err := wiki.GetNodes()
		if err != nil {
			panic(err)
		}

		for _, node := range nodes {
			tasks := node.GetTasks()
			for _, task := range tasks {
				if task.IsInProgress() {
					if cmd.Flag("elapsed").Value != nil {
						start := task.GetLastWorkSession().Start
						elapsed := time.Since(start).Round(time.Second)
						fmt.Printf("%s - %s (%s)\n", node.GetName(), task.Text, elapsed)
					} else {
						fmt.Printf("%s - %s\n", node.GetName(), task.Text)
					}
					return
				}
			}
		}
	},
}

func init() {
	cmd := &cobra.Command{Use: "task"}

	currentTaskCommand.Flags().BoolP("elapsed", "e", false, "include elapsed time")
	cmd.AddCommand(currentTaskCommand)

	rootCmd.AddCommand(cmd)
}
