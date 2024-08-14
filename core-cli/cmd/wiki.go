package cmd

import (
	wikivfs "core/vfs/wiki-vfs"
	"fmt"
	"log"
	"path/filepath"

	local_wiki "github.com/3rd/core/core-lib/wiki/local"
	"github.com/radovskyb/watcher"
	"github.com/spf13/cobra"
)

var wikiListCommand = &cobra.Command{
	Use:   "ls",
	Short: "list wiki nodes",
	Run: func(cmd *cobra.Command, args []string) {
		root := env.WIKI_ROOT
		if len(root) == 0 {
			panic("WIKI_ROOT not set")
		}

		isDebug, err := cmd.Flags().GetBool("debug")
		if err != nil {
			panic(err)
		}

		typeFilter, err := cmd.Flags().GetString("type")
		if err != nil {
			panic(err)
		}

		// regular
		if !isDebug {
			wiki, err := local_wiki.NewLocalWiki(local_wiki.LocalWikiConfig{
				Root:  root,
				Parse: "meta",
			})
			if err != nil {
				panic(err)
			}
			nodes, _ := wiki.GetNodes()

			for _, node := range nodes {
				meta := node.GetMeta()
				if meta != nil && typeFilter == "" || meta["type"] == typeFilter {
					fmt.Printf("%s\n", node.GetID())
				}
			}
		}

		// debug
		if isDebug {
			wiki, err := local_wiki.NewLocalWiki(local_wiki.LocalWikiConfig{
				Root:  root,
				Parse: "full",
			})
			if err != nil {
				panic(err)
			}
			nodes, _ := wiki.GetNodes()
			for _, node := range nodes {
				meta := node.GetMeta()
				if meta != nil && typeFilter == "" || meta["type"] == typeFilter {
					fmt.Printf("%s %s\n", node.GetID(), node.ParseDuration)
				}
			}
		}
	},
}

var wikiResolveCommand = &cobra.Command{
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

		wiki, err := local_wiki.NewLocalWiki(local_wiki.LocalWikiConfig{
			Root:  root,
			Parse: "meta",
		})
		if err != nil {
			panic(err)
		}

		target := args[0]

		nodes, err := wiki.GetNodes()
		if err != nil {
			panic(err)
		}
		for _, node := range nodes {
			if node.GetName() == target {
				fmt.Print(node.GetPath())
				return
			}
		}

		if !isStrict {
			unsortedPath := filepath.Join(env.WIKI_ROOT, "unsorted", target)
			fmt.Print(unsortedPath)
		}
	},
}

var wikiMountCommand = &cobra.Command{
	Use:   "mount",
	Short: "mount wiki vfs",
	Run: func(cmd *cobra.Command, args []string) {
		root := env.WIKI_ROOT
		if len(root) == 0 {
			panic("WIKI_ROOT not set")
		}

		mountPoint, err := cmd.Flags().GetString("mount")
		if err != nil {
			panic(err)
		}

		// get wiki
		wikiInstance, err := local_wiki.NewLocalWiki(local_wiki.LocalWikiConfig{
			Root:  root,
			Parse: "none",
		})
		if err != nil {
			panic(err)
		}

		// setup watcher
		w := watcher.New()
		w.FilterOps(watcher.Create, watcher.Move, watcher.Remove, watcher.Write)

		go func() {
			for {
				select {
				case <-w.Event:
					wikiInstance.Reload()
				case err := <-w.Error:
					log.Fatalln(err)
				case <-w.Closed:
					return
				}
			}
		}()
		if err := w.AddRecursive(root); err != nil {
			log.Fatalln(err)
		}

		// mount
		vfs, err := wikivfs.NewWikiVFS(wikiInstance, root, mountPoint)
		if err != nil {
			panic(err)
		}
		defer vfs.Close()
		err = vfs.Mount()
		if err != nil {
			panic(err)
		}
	},
}

func init() {
	cmd := &cobra.Command{Use: "wiki"}

	wikiListCommand.Flags().Bool("debug", false, "debug parsing time for each nodes")
	wikiListCommand.Flags().String("type", "", "filter nodes by type")
	cmd.AddCommand(wikiListCommand)

	wikiResolveCommand.Flags().Bool("strict", false, "will not return the default would-be path for if the node is not found")
	cmd.AddCommand(wikiResolveCommand)

	wikiMountCommand.Flags().String("mount", "/tmp/wiki", "mount point")
	cmd.AddCommand(wikiMountCommand)

	rootCmd.AddCommand(cmd)
}
