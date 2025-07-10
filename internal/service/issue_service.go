package service

import (
    "errors"
    "time"

    "backed-assignment-aro/internal/model"
    "backed-assignment-aro/internal/repository"
)

// IssueService contains business logic around issues.
type IssueService struct {
    repo *repository.MemoryRepository
}

func NewIssueService(r *repository.MemoryRepository) *IssueService {
    return &IssueService{repo: r}
}

// CreateIssue validates and stores a new issue.
func (s *IssueService) CreateIssue(title, desc string, userID *uint) (*model.Issue, error) {
    if title == "" {
        return nil, errors.New("title is required")
    }

    var usr *model.User
    var status model.IssueStatus
    if userID != nil {
        u, ok := s.repo.GetUser(*userID)
        if !ok {
            return nil, errors.New("user not found")
        }
        usr = u
        status = model.StatusInProgress
    } else {
        status = model.StatusPending
    }

    now := time.Now().UTC()
    iss := &model.Issue{
        Title:       title,
        Description: desc,
        Status:      status,
        User:        usr,
        CreatedAt:   now,
        UpdatedAt:   now,
    }
    return s.repo.CreateIssue(iss), nil
}

// UpdateIssue applies update rules.
func (s *IssueService) UpdateIssue(id uint, title, desc *string, status *model.IssueStatus, userID *uint) (*model.Issue, error) {
    iss, ok := s.repo.GetIssue(id)
    if !ok {
        return nil, errors.New("issue not found")
    }

    // Prevent updates on completed/cancelled
    if iss.Status == model.StatusCompleted || iss.Status == model.StatusCancelled {
        return nil, errors.New("cannot update completed or cancelled issue")
    }

    // Handle status/user rules
    if userID != nil {
        if *userID == 0 { // remove user
            iss.User = nil
            iss.Status = model.StatusPending
        } else {
            u, ok := s.repo.GetUser(*userID)
            if !ok {
                return nil, errors.New("user not found")
            }
            iss.User = u
            if status == nil {
                // auto change depending on previous state
                if iss.Status == model.StatusPending {
                    iss.Status = model.StatusInProgress
                }
            }
        }
    }

    if status != nil {
        // validate state transitions
        if (iss.User == nil) && (*status != model.StatusPending && *status != model.StatusCancelled) {
            return nil, errors.New("cannot set status without assignee")
        }
        iss.Status = *status
    }

    if title != nil {
        iss.Title = *title
    }
    if desc != nil {
        iss.Description = *desc
    }

    iss.UpdatedAt = time.Now().UTC()
    if err := s.repo.UpdateIssue(iss); err != nil {
        return nil, err
    }
    return iss, nil
}

func (s *IssueService) ListIssues(status *model.IssueStatus) []*model.Issue {
    return s.repo.ListIssues(status)
}

func (s *IssueService) GetIssue(id uint) (*model.Issue, error) {
    iss, ok := s.repo.GetIssue(id)
    if !ok {
        return nil, errors.New("issue not found")
    }
    return iss, nil
}
