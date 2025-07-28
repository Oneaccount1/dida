package entities

import "time"

// Project 表示TickTick项目实体
type Project struct {
	ID        string
	Name      string
	Color     string
	ViewMode  ViewMode
	Closed    bool
	Kind      ProjectKind
	CreatedAt time.Time
	UpdatedAt time.Time
}

// ViewMode 项目视图模式
type ViewMode string

const (
	ViewModeList     ViewMode = "list"
	ViewModeBoard    ViewMode = "board"
	ViewModeCalendar ViewMode = "calendar"
	ViewModeTimeline ViewMode = "timeline"
)

// ProjectKind 项目类型
type ProjectKind string

const (
	ProjectKindTask     ProjectKind = "task"
	ProjectKindNote     ProjectKind = "note"
	ProjectKindHabit    ProjectKind = "habit"
)

// Business methods for Project entity

// IsActive 检查项目是否活跃
func (p *Project) IsActive() bool {
	return !p.Closed
}

// Close 关闭项目
func (p *Project) Close() {
	p.Closed = true
	p.UpdatedAt = time.Now()
}

// Reopen 重新打开项目
func (p *Project) Reopen() {
	p.Closed = false
	p.UpdatedAt = time.Now()
}

// SetViewMode 设置视图模式
func (p *Project) SetViewMode(mode ViewMode) {
	p.ViewMode = mode
	p.UpdatedAt = time.Now()
}

// SetColor 设置项目颜色
func (p *Project) SetColor(color string) {
	p.Color = color
	p.UpdatedAt = time.Now()
}