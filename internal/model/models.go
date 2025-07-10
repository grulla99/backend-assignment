package model

import "time"

type User struct {
	ID   uint   `json:"id"`
	Name string `json:"name"`
}

type IssueStatus string

const (
	StatusPending    IssueStatus = "PENDING"
	StatusInProgress IssueStatus = "IN_PROGRESS"
	StatusCompleted  IssueStatus = "COMPLETED"
	StatusCancelled  IssueStatus = "CANCELLED"
)

type Issue struct {
	ID          uint        `json:"id"`
	Title       string      `json:"title"`
	Description string      `json:"description,omitempty"`
	Status      IssueStatus `json:"status"`
	User        *User       `json:"user,omitempty"`
	CreatedAt   time.Time   `json:"createdAt"`
	UpdatedAt   time.Time   `json:"updatedAt"`
}
