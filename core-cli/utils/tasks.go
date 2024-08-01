package utils

import "github.com/3rd/core/core-lib/wiki"

func ComputeTaskReward(task *wiki.Task) int {
	points := task.Priority
	if points == 0 {
		points = 1
	}
	if task.Schedule != nil && task.Schedule.Repeat != "" {
		points = 10
	}
	points = points * 10
	return int(points)
}
