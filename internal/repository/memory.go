package repository

import (
	"errors"
	"sync"

	"backed-assignment-aro/internal/model"
)

// 메모리 저장소 구조체
type MemoryRepository struct {
	mu     sync.RWMutex          // 동시성 제어 뮤텍스
	users  map[uint]*model.User  // 사용자 저장소
	issues map[uint]*model.Issue // 이슈 저장소
	nextID uint                  // 다음 이슈 id
}

// 새 메모리 저장소 생성 및 초기 사용자 설정
func NewMemoryRepository() *MemoryRepository {
	r := &MemoryRepository{
		users:  make(map[uint]*model.User),
		issues: make(map[uint]*model.Issue),
		nextID: 1,
	}

	// 초기 사용자 3명
	r.users[1] = &model.User{ID: 1, Name: "김개발"}
	r.users[2] = &model.User{ID: 2, Name: "이디자인"}
	r.users[3] = &model.User{ID: 3, Name: "박기획"}
	return r
}

// 사용자 조회
func (r *MemoryRepository) GetUser(id uint) (*model.User, bool) {
	r.mu.RLock()         // 읽기 잠금
	defer r.mu.RUnlock() // 한수 종료 시 잠금 해제

	u, ok := r.users[id]
	return u, ok
}

// 이슈 목록 조회(상태 필터링)
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

// 이슈 상세 조회
func (r *MemoryRepository) GetIssue(id uint) (*model.Issue, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	iss, ok := r.issues[id]
	return iss, ok
}

// 이슈 생성
func (r *MemoryRepository) CreateIssue(iss *model.Issue) *model.Issue {
	r.mu.Lock() // 쓰기 잠금
	defer r.mu.Unlock()

	iss.ID = r.nextID
	r.nextID++
	r.issues[iss.ID] = iss
	return iss
}

// 이슈 수정
func (r *MemoryRepository) UpdateIssue(iss *model.Issue) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, ok := r.issues[iss.ID]; !ok {
		return errors.New("issue not found")
	}
	r.issues[iss.ID] = iss
	return nil
}
