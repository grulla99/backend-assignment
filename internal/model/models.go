package model

import "time"

// 사용자 정보
type User struct {
	ID   uint   `json:"id"`
	Name string `json:"name"`
}

// 이슈 상태 타입
type IssueStatus string

const (
	StatusPending    IssueStatus = "PENDING"
	StatusInProgress IssueStatus = "IN_PROGRESS"
	StatusCompleted  IssueStatus = "COMPLETED"
	StatusCancelled  IssueStatus = "CANCELLED"
)

// 이슈 정보
type Issue struct {
	ID          uint        `json:"id"`
	Title       string      `json:"title"`
	Description string      `json:"description,omitempty"`
	Status      IssueStatus `json:"status"`
	User        *User       `json:"user,omitempty"`
	CreatedAt   time.Time   `json:"createdAt"`
	UpdatedAt   time.Time   `json:"updatedAt"`
}
