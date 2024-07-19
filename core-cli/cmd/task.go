package cmd

import (
	taskinteractive "core/ui/task_interactive"
	"fmt"
	"sort"
	"time"

	"github.com/3rd/core/core-lib/wiki"
	localWiki "github.com/3rd/core/core-lib/wiki/local"
	"github.com/spf13/cobra"
)

var taskCurrentCommand = &cobra.Command{
	Use:   "current",
	Short: "list the currently in-progress task (first only)",
	Run: func(cmd *cobra.Command, args []string) {
		root := env.WIKI_ROOT
		if len(root) == 0 {
			panic("WIKI_ROOT not set")
		}

		wiki, err := localWiki.NewLocalWiki(localWiki.LocalWikiConfig{
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

var taskInteractiveCommand = &cobra.Command{
	Use:   "interactive",
	Short: "enter interactive task mode",
	Run: func(cmd *cobra.Command, args []string) {
		root := env.TASK_ROOT
		if len(root) == 0 {
			panic("TASK_ROOT not set")
		}

		wikiInstance, err := localWiki.NewLocalWiki(localWiki.LocalWikiConfig{
			Root:            root,
			Parse:           "full",
			SkipInitialLoad: true,
		})
		if err != nil {
			panic(err)
		}

		loadTasks := func() taskinteractive.GetTasksResult {
			wikiInstance.Reload()

			nodes, err := wikiInstance.GetNodes()
			if err != nil {
				panic(err)
			}

			// load tasks
			recentlyDoneLimit, _ := time.ParseDuration("24h")
			tasks := []*wiki.Task{}
			longestProjectLength := 0
			now := time.Now()

			for _, node := range nodes {
				nodeTasks := node.GetTasks()
				hasAddedTaskForNode := false

				for _, task := range nodeTasks {
					if task.Status == wiki.TASK_STATUS_ACTIVE || task.IsInProgress() {
						tasks = append(tasks, task)
						hasAddedTaskForNode = true
					}
					if task.Status == wiki.TASK_STATUS_DONE {
						last_session := task.GetLastWorkSession()
						if last_session != nil && now.Sub(last_session.Start) < recentlyDoneLimit {
							// todayStart := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.Local)
							// if last_session != nil && last_session.Start.After(todayStart) {
							tasks = append(tasks, task)
							hasAddedTaskForNode = true
						}
					}
				}

				if hasAddedTaskForNode {
					projectLength := len(node.GetName())
					if projectLength > longestProjectLength {
						longestProjectLength = projectLength
					}
				}
			}

			// custom sort
			sort.Slice(tasks, func(i, j int) bool {
				a := tasks[i]
				b := tasks[j]

				// top: sticky
				if b.Priority >= 100 {
					return false
				}
				if a.Priority >= 100 {
					return true
				}

				// middle: done
				if b.Status == wiki.TASK_STATUS_DONE && a.Status != wiki.TASK_STATUS_DONE {
					return false
				}
				if a.Status == wiki.TASK_STATUS_DONE && b.Status != wiki.TASK_STATUS_DONE {
					return true
				}

				// by priority
				if a.Priority != b.Priority {
					return b.Priority < a.Priority
				}

				// by location
				aNode := a.Node.(*localWiki.LocalNode)
				bNode := b.Node.(*localWiki.LocalNode)
				if aNode.GetPath() != bNode.GetPath() {
					return a.LineNumber < b.LineNumber
				}
				return a.Text < b.Text
			})

			return taskinteractive.GetTasksResult{
				Tasks:                tasks,
				LongestProjectLength: longestProjectLength,
			}
		}

		providers := taskinteractive.Providers{
			GetTasks: loadTasks,
		}
		taskinteractive.Run(providers)
	},
}

func init() {
	cmd := &cobra.Command{Use: "task"}

	taskCurrentCommand.Flags().BoolP("elapsed", "e", false, "include elapsed time")
	cmd.AddCommand(taskCurrentCommand)

	cmd.AddCommand(taskInteractiveCommand)

	rootCmd.AddCommand(cmd)
}
