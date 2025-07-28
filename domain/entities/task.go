package entities

import (
	"time"
)

// Task 表示TickTick任务实体
type Task struct {
	ID            string
	ProjectID     string
	Title         string
	Content       string
	Description   string
	IsAllDay      bool
	StartDate     *time.Time
	DueDate       *time.Time
	TimeZone      string
	Reminders     []string
	RepeatFlag    string
	Priority      Priority
	Status        TaskStatus
	CompletedTime *time.Time
	SortOrder     int
	Items         []TaskItem
	CreatedAt     time.Time
	UpdatedAt     time.Time
}

// TaskItem 表示子任务
type TaskItem struct {
	ID            string
	Status        TaskStatus
	Title         string
	SortOrder     int
	StartDate     *time.Time
	IsAllDay      bool
	TimeZone      string
	CompletedTime *time.Time
}

// Priority 任务优先级
type Priority int

const (
	PriorityNone Priority = iota
	PriorityLow
	PriorityMedium
	PriorityHigh
)

// TaskStatus 任务状态
type TaskStatus int

const (
	TaskStatusIncomplete TaskStatus = iota
	TaskStatusCompleted
)

// Business methods for Task entity

// IsCompleted 检查任务是否已完成
func (t *Task) IsCompleted() bool {
	return t.Status == TaskStatusCompleted
}

// IsOverdue 检查任务是否已过期
func (t *Task) IsOverdue() bool {
	if t.DueDate == nil {
		return false
	}
	return time.Now().After(*t.DueDate) && !t.IsCompleted()
}

// Complete 完成任务
func (t *Task) Complete() {
	now := time.Now()
	t.Status = TaskStatusCompleted
	t.CompletedTime = &now
	t.UpdatedAt = now
}

// SetPriority 设置任务优先级
func (t *Task) SetPriority(priority Priority) {
	t.Priority = priority
	t.UpdatedAt = time.Now()
}

// AddSubTask 添加子任务
func (t *Task) AddSubTask(title string) *TaskItem {
	item := TaskItem{
		Title:     title,
		Status:    TaskStatusIncomplete,
		SortOrder: len(t.Items),
	}
	t.Items = append(t.Items, item)
	t.UpdatedAt = time.Now()
	return &item
}