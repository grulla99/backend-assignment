package service

import (
	"errors"
	"time"

	"backed-assignment-aro/internal/model"
	"backed-assignment-aro/internal/repository"
)

// 이슈 서비스 구조체
type IssueService struct {
	repo *repository.MemoryRepository
}

// 새 이슈 서비스 생성
func NewIssueService(r *repository.MemoryRepository) *IssueService {
	return &IssueService{repo: r}
}

// 이슈 생성 로직
func (s *IssueService) CreateIssue(title, desc string, userID *uint) (*model.Issue, error) {
	// 제목 필수 검증
	if title == "" {
		return nil, errors.New("title is required")
	}

	var usr *model.User
	var status model.IssueStatus

	// 사용자 지정 시 상태 설정
	if userID != nil {
		u, ok := s.repo.GetUser(*userID)
		if !ok {
			return nil, errors.New("user not found")
		}
		usr = u
		status = model.StatusInProgress // 지정하면 진행 중
	} else {
		status = model.StatusPending // 미지정 시 대기 중
	}

	// 새 이슈 생성
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

// 이슈 수정
func (s *IssueService) UpdateIssue(id uint, title, desc *string, status *model.IssueStatus, userID *uint) (*model.Issue, error) {
	// 기존 이슈 조회
	iss, ok := s.repo.GetIssue(id)
	if !ok {
		return nil, errors.New("issue not found")
	}

	// 완료,취소 된 이슈는 수정 불가
	if iss.Status == model.StatusCompleted || iss.Status == model.StatusCancelled {
		return nil, errors.New("cannot update completed or cancelled issue")
	}

	// 사용자 변경 로직
	if userID != nil {
		if *userID == 0 { // 사용자 제거
			iss.User = nil
			iss.Status = model.StatusPending
		} else { // 새 사용자 할당
			u, ok := s.repo.GetUser(*userID)
			if !ok {
				return nil, errors.New("user not found")
			}
			iss.User = u
			if status == nil {
				// 만약 상태가 지정되지 않고 현재 상태가 대기중이면 진행 중으로 변경
				if iss.Status == model.StatusPending {
					iss.Status = model.StatusInProgress
				}
			}
		}
	}

	// 상태 유효성 검증
	if status != nil {
		// 담당자 없이 진행중/완료 상태로 변경 불가
		if (iss.User == nil) && (*status != model.StatusPending && *status != model.StatusCancelled) {
			return nil, errors.New("cannot set status without assignee")
		}
		iss.Status = *status
	}

	// 필드 업데이트
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
