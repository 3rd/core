package cmd

import (
	"core/utils"
	"fmt"
	"path/filepath"

	wiki "github.com/3rd/core/core-lib/wiki/local"
	"github.com/spf13/cobra"
)

var env = utils.GetEnv()

var listCommand = &cobra.Command{
	Use:   "ls",
	Short: "list wiki nodes",
	Run: func(cmd *cobra.Command, args []string) {
		root := env.WIKI_ROOT
		if len(root) == 0 {
			panic("WIKI_ROOT not set")
		}

		wiki, err := wiki.NewLocalWiki(wiki.LocalWikiConfig{
			Root:  root,
			Parse: "meta",
		})
		if err != nil {
			panic(err)
		}

		nodes, _ := wiki.GetNodes()

		for _, node := range nodes {
			fmt.Printf("%v\n", node.GetID())
		}
	},
}

var resolveCommand = &cobra.Command{
	Use:   "resolve <node>",
	Short: "show node file path",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			panic("No node specified")
		}

		isStrict, err := cmd.Flags().GetBool("strict")
		if err != nil {
			panic(err)
		}

		root := env.WIKI_ROOT
		if len(root) == 0 {
			panic("WIKI_ROOT not set")
		}
		wiki, err := wiki.NewLocalWiki(wiki.LocalWikiConfig{
			Root:  root,
			Parse: "meta",
		})
		if err != nil {
			panic(err)
		}

		target := args[0]

		if node, _ := wiki.GetNode(target); node != nil {
			fmt.Print(node.GetPath())
		} else {
			if !isStrict {
				unsortedPath := filepath.Join(env.WIKI_ROOT, "unsorted", target)
				fmt.Print(unsortedPath)
			}
		}
	},
}

func init() {
	cmd := &cobra.Command{Use: "wiki"}

	cmd.AddCommand(listCommand)
	cmd.AddCommand(resolveCommand)
	resolveCommand.Flags().Bool("strict", false, "Will not return the default would-be path for if the node is not found")

	rootCmd.AddCommand(cmd)
}
