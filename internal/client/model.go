package client

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
	ID       string `json:"id,omitempty"`
	Name     string `json:"name"`
	Color    string `json:"color,omitempty"`
	ViewMode string `json:"viewMode,omitempty"`
	Closed   bool   `json:"closed,omitempty"`
	Kind     string `json:"kind,omitempty"`
}
