package cmd

import (
	taskinteractive "core/ui/task_interactive"
	"fmt"
	"sort"
	"strings"
	"time"

	"github.com/3rd/core/core-lib/wiki"
	localWiki "github.com/3rd/core/core-lib/wiki/local"
	"github.com/spf13/cobra"
)

func preprocessScheduleRepeat(input string) string {
	switch input {
	case "day":
		return "mon,tue,wed,thu,fri,sat,sun"
	case "week":
		return "mon,tue,wed,thu,fri,sat,sun"
	case "workday":
		return "mon,tue,wed,thu,fri"
	}
	return input
}

func formatRepeatDayStr(input string) string {
	switch input {
	case "mon":
		return "Monday"
	case "tue":
		return "Tuesday"
	case "wed":
		return "Wednesday"
	case "thu":
		return "Thursday"
	case "fri":
		return "Friday"
	case "sat":
		return "Saturday"
	case "sun":
		return "Sunday"
	}
	return input
}

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
			meta := node.GetMeta()
			nodeType, ok := meta["type"]
			if !ok || nodeType != "project" {
				continue
			}

			tasks := node.GetTasks()
			for _, task := range tasks {
				if task.IsInProgress() {
					if cmd.Flag("elapsed").Value != nil {
						start := task.GetLastSession().Start
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
			activeTasks := []*wiki.Task{}
			tasks := []*wiki.Task{}
			longestActiveProjectLength := 0

			recentlyDoneOffset, _ := time.ParseDuration("24h")
			now := time.Now()
			startOfDay := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.Local)
			if recentlyDoneOffset > 0 {
				startOfDay = startOfDay.Add(-recentlyDoneOffset)
			}
			endOfDay := time.Date(now.Year(), now.Month(), now.Day()+1, 0, 0, 0, 0, time.Local)

			for _, node := range nodes {
				meta := node.GetMeta()
				nodeType, ok := meta["type"]
				if !ok || nodeType != "project" {
					continue
				}

				nodeTasks := node.GetTasks()
				hasAddedActiveTaskForNode := false

				for _, task := range nodeTasks {
					var taskToAdd *wiki.Task

					// skip cancelled
					if task.Status == wiki.TASK_STATUS_CANCELLED {
						continue
					}

					// active
					if task.Status == wiki.TASK_STATUS_ACTIVE {
						taskToAdd = task
					}

					// done today
					if taskToAdd == nil && task.Status == wiki.TASK_STATUS_DONE {
						for _, session := range task.Sessions {
							if session.Start.After(startOfDay) && session.Start.Before(endOfDay) {
								taskToAdd = task
								break
							}
						}
					}

					// scheduled for today
					if taskToAdd == nil && task.Status == wiki.TASK_STATUS_DEFAULT && task.Schedule != nil {
						if task.Schedule.Start.After(startOfDay) && task.Schedule.Start.Before(endOfDay) {
							taskToAdd = task
						}
					}

					// scheduled in the past, not completed, without recurrence
					if taskToAdd == nil && task.Status != wiki.TASK_STATUS_DONE && task.Schedule != nil && task.Schedule.Repeat == "" {
						if task.Schedule.Start.Before(startOfDay) {
							taskToAdd = task
						}
					}

					// with a recurring schedule that is due today
					if taskToAdd == nil && task.Status != wiki.TASK_STATUS_DONE && task.Schedule != nil && task.Schedule.Repeat != "" {
						isShortNotation := false

						// @daily
						if task.Schedule.Repeat == "daily" {
							isShortNotation = true
							taskToAdd = task
						}
						// @weekly
						if task.Schedule.Repeat == "weekly" {
							isShortNotation = true
							scheduledDay := task.Schedule.Start.Weekday()
							if scheduledDay == now.Weekday() {
								taskToAdd = task
							}
						}
						// @monthly
						if task.Schedule.Repeat == "monthly" {
							isShortNotation = true
							scheduledDay := task.Schedule.Start.Day()
							if scheduledDay == now.Day() {
								taskToAdd = task
							}
						}

						// mon, tue, wed, thu, fri, sat, sun
						if !isShortNotation {
							parts := strings.Split(preprocessScheduleRepeat(task.Schedule.Repeat), ",")
							for _, part := range parts {
								if formatRepeatDayStr(part) == now.Weekday().String() {
									taskToAdd = task
								}
							}
						}
					}

					// patch completion
					if task.HasCompletionForDate(now) {
						task.Status = wiki.TASK_STATUS_DONE
					}

					// add task
					tasks = append(tasks, task)
					if taskToAdd != nil {
						activeTasks = append(activeTasks, taskToAdd)
						hasAddedActiveTaskForNode = true
					}
				}

				projectLength := len(node.GetName())
				if hasAddedActiveTaskForNode {
					if projectLength > longestActiveProjectLength {
						longestActiveProjectLength = projectLength
					}
				}
			}

			// custom sort
			sort.Slice(activeTasks, func(i, j int) bool {
				a := activeTasks[i]
				b := activeTasks[j]

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

				// by schedule
				if a.Schedule != nil && b.Schedule != nil {
					aStart := a.Schedule.Start.Hour()*60 + a.Schedule.Start.Minute()
					bStart := b.Schedule.Start.Hour()*60 + b.Schedule.Start.Minute()
					if aStart != bStart {
						if aStart == 0 && bStart != 0 {
							return false
						}
						if aStart != 0 && bStart == 0 {
							return true
						}
						return aStart < bStart
					}
				}
				if a.Schedule != nil && b.Schedule == nil {
					return true
				}
				if a.Schedule == nil && b.Schedule != nil {
					return false
				}

				// by location
				aNode := a.Node.(*localWiki.LocalNode)
				bNode := b.Node.(*localWiki.LocalNode)
				if aNode.GetPath() == bNode.GetPath() {
					return a.LineNumber < b.LineNumber
				}
				return a.Node.GetName() < b.Node.GetName()
			})

			return taskinteractive.GetTasksResult{
				Tasks:                      tasks,
				ActiveTasks:                activeTasks,
				LongestActiveProjectLength: longestActiveProjectLength,
			}
		}

		getRoot := func() string {
			return root
		}

		providers := taskinteractive.Providers{
			GetTasks: loadTasks,
			GetRoot:  getRoot,
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
