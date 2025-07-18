package client

import "time"

// Task 表示TickTick任务
type Task struct {
	ID            string     `json:"id,omitempty"`
	ProjectID     string     `json:"projectId"`
	Title         string     `json:"title"`
	Content       string     `json:"content,omitempty"`
	Desc          string     `json:"desc,omitempty"`
	IsAllDay      bool       `json:"isAllDay,omitempty"`
	StartDate     string     `json:"startDate,omitempty"`
	DueDate       string     `json:"dueDate,omitempty"`
	TimeZone      string     `json:"timeZone,omitempty"`
	Reminders     []string   `json:"reminders,omitempty"`
	RepeatFlag    string     `json:"repeatFlag,omitempty"`
	Priority      int        `json:"priority,omitempty"`
	Status        int        `json:"status,omitempty"`
	CompletedTime string     `json:"completedTime,omitempty"`
	SortOrder     int        `json:"sortOrder,omitempty"`
	Items         []TaskItem `json:"items,omitempty"`
}

// TaskItem 表示子任务
type TaskItem struct {
	ID            string `json:"id,omitempty"`
	Status        int    `json:"status"`
	Title         string `json:"title"`
	SortOrder     int    `json:"sortOrder,omitempty"`
	StartDate     string `json:"startDate,omitempty"`
	IsAllDay      bool   `json:"isAllDay,omitempty"`
	TimeZone      string `json:"timeZone,omitempty"`
	CompletedTime string `json:"completedTime,omitempty"`
}

// Project 表示TickTick项目
type Project struct {
	ID        string `json:"id,omitempty"`
	Name      string `json:"name"`
	Color     string `json:"color,omitempty"`
	ViewMode  string `json:"viewMode,omitempty"`
	SortOrder int64  `json:"sortOrder,omitempty"`
	Kind      string `json:"kind,omitempty"`
}

// ChecklistItem 表示一个待办事项清单中的子任务
type ChecklistItem struct {
	ID            string    `json:"id"`            // 子任务标识符
	Title         string    `json:"title"`         // 子任务标题
	Status        int32     `json:"status"`        // 子任务完成状态: 0=Normal, 1=Completed
	CompletedTime time.Time `json:"completedTime"` // 子任务完成时间
	IsAllDay      bool      `json:"isAllDay"`      // 是否全天任务
	SortOrder     int64     `json:"sortOrder"`     // 子任务排序顺序
	StartDate     time.Time `json:"startDate"`     // 子任务开始时间
	TimeZone      string    `json:"timeZone"`      // 子任务时区
}

type Column struct {
	ID        string `json:"id"`
	ProjectID string `json:"projectID"`
	Name      string `json:"Name"`
	SortOrder int64  `json:"sortOrder"`
}

type ProjectData struct {
	Project Project
	Tasks   []Task
	Columns []Column
}
