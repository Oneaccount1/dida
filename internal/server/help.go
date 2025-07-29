package server

import (
	"dida/internal/client"
	"fmt"
)

func FormatProject(project client.Project) string {
	formatted := fmt.Sprintf("Name: %s\n", project.Name)
	formatted += fmt.Sprintf("ID: %s\n", project.ID)

	if project.Color != "" {
		formatted += fmt.Sprintf("Color: %s\n", project.Color)
	}
	if project.ViewMode != "" {
		formatted += fmt.Sprintf("View Mode: %s\n", project.ViewMode)
	}
	if project.Kind != "" {
		formatted += fmt.Sprintf("Kind: %s\n", project.Kind)
	}
	return formatted
}

// FormatTask 将任务对象格式化为可读字符串
func FormatTask(task client.Task) string {

	formatted := fmt.Sprintf("ID: %s\n", task.ID)
	formatted += fmt.Sprintf("Title: %s\n", task.Title)
	formatted += fmt.Sprintf("Project ID: %s\n", task.ProjectID)

	if task.StartDate != "" {
		formatted += fmt.Sprintf("Start Date: %s\n", task.StartDate)
	}

	if task.DueDate != "" {
		formatted += fmt.Sprintf("Due Date: %s\n", task.DueDate)
	}

	priorityMap := map[int]string{
		0: "None",
		1: "Low",
		3: "Medium",
		5: "High",
	}
	formatted += fmt.Sprintf("Priority: %s\n", priorityMap[task.Priority])

	status := "Active"
	if task.Status == 2 {
		status = "Completed"
	}
	formatted += fmt.Sprintf("Status: %s\n", status)

	if task.Content != "" {
		formatted += fmt.Sprintf("\nContent:\n%s\n", task.Content)
	}

	if len(task.Items) > 0 {
		formatted += fmt.Sprintf("\nSubtasks (%d):\n", len(task.Items))
		for i, item := range task.Items {
			statusMark := "□"
			if item.Status == 1 {
				statusMark = "✓"
			}
			formatted += fmt.Sprintf("%d. [%s] %s\n", i+1, statusMark, item.Title)
		}
	}

	return formatted
}
