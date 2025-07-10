package repository

import (
    "errors"
    "sync"

    "backed-assignment-aro/internal/model"
)

// MemoryRepository provides an in-memory implementation of repositories for users and issues.
// It is safe for concurrent use.
type MemoryRepository struct {
    mu     sync.RWMutex
    users  map[uint]*model.User
    issues map[uint]*model.Issue
    nextID uint
}

// NewMemoryRepository returns an initialized repository with seed users.
func NewMemoryRepository() *MemoryRepository {
    r := &MemoryRepository{
        users:  make(map[uint]*model.User),
        issues: make(map[uint]*model.Issue),
        nextID: 1,
    }
    // Seed default users as required.
    r.users[1] = &model.User{ID: 1, Name: "김개발"}
    r.users[2] = &model.User{ID: 2, Name: "이디자인"}
    r.users[3] = &model.User{ID: 3, Name: "박기획"}
    return r
}

// GetUser returns a user by id.
func (r *MemoryRepository) GetUser(id uint) (*model.User, bool) {
    r.mu.RLock()
    defer r.mu.RUnlock()
    u, ok := r.users[id]
    return u, ok
}

// ListIssues optionally filters by status.
func (r *MemoryRepository) ListIssues(status *model.IssueStatus) []*model.Issue {
    r.mu.RLock()
    defer r.mu.RUnlock()
    var res []*model.Issue
    for _, iss := range r.issues {
        if status == nil || iss.Status == *status {
            res = append(res, iss)
        }
    }
    return res
}

// GetIssue retrieves an issue by id.
func (r *MemoryRepository) GetIssue(id uint) (*model.Issue, bool) {
    r.mu.RLock()
    defer r.mu.RUnlock()
    iss, ok := r.issues[id]
    return iss, ok
}

// CreateIssue stores a new issue and returns it.
func (r *MemoryRepository) CreateIssue(iss *model.Issue) *model.Issue {
    r.mu.Lock()
    defer r.mu.Unlock()
    iss.ID = r.nextID
    r.nextID++
    r.issues[iss.ID] = iss
    return iss
}

// UpdateIssue updates an existing issue in place.
func (r *MemoryRepository) UpdateIssue(iss *model.Issue) error {
    r.mu.Lock()
    defer r.mu.Unlock()
    if _, ok := r.issues[iss.ID]; !ok {
        return errors.New("issue not found")
    }
    r.issues[iss.ID] = iss
    return nil
}
