package handler

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5"

	"backed-assignment-aro/internal/model"
	"backed-assignment-aro/internal/service"
)

// HTTP 핸들러 구조체
type HTTPHandler struct {
	svc *service.IssueService
}

// 새 HTTP 핸들러 생성
func NewHTTPHandler(s *service.IssueService) *HTTPHandler {
	return &HTTPHandler{svc: s}
}

// 라우트 등록
func (h *HTTPHandler) Register(r chi.Router) {
	r.Post("/issue", h.createIssue)       // 생성
	r.Get("/issues", h.listIssues)        // 목록 조회
	r.Get("/issue/{id}", h.getIssue)      // 상세 조회
	r.Patch("/issue/{id}", h.updateIssue) // 수정
}

// 에러 응답 처리
func respondErr(w http.ResponseWriter, code int, msg string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	_ = json.NewEncoder(w).Encode(map[string]interface{}{
		"error": msg,
		"code":  code,
	})
}

// 이슈 생성 핸들러
func (h *HTTPHandler) createIssue(w http.ResponseWriter, r *http.Request) {
	var body struct {
		Title       string `json:"title"`
		Description string `json:"description"`
		UserID      *uint  `json:"userId"`
	}

	// 요청 바디 파싱
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		respondErr(w, http.StatusBadRequest, "invalid json")
		return
	}

	// 서비스 호출
	iss, err := h.svc.CreateIssue(body.Title, body.Description, body.UserID)
	if err != nil {
		respondErr(w, http.StatusBadRequest, err.Error())
		return
	}

	// 성공 응답
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(iss)
}

type IssueResponse struct {
	ID          uint      `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description,omitempty"`
	Status      string    `json:"status"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
}

type IssuesListResponse struct {
	Issues []IssueResponse `json:"issues"`
}

func (h *HTTPHandler) listIssues(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query().Get("status")
	var status *model.IssueStatus
	if q != "" {
		s := model.IssueStatus(q)
		status = &s
	}

	issues := h.svc.ListIssues(status)

	var response IssuesListResponse
	for _, issue := range issues {
		response.Issues = append(response.Issues, IssueResponse{
			ID:          issue.ID,
			Title:       issue.Title,
			Description: issue.Description,
			Status:      string(issue.Status),
			CreatedAt:   issue.CreatedAt,
			UpdatedAt:   issue.UpdatedAt,
		})
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(response)
}

func (h *HTTPHandler) getIssue(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		respondErr(w, http.StatusBadRequest, "invalid id")
		return
	}
	iss, err := h.svc.GetIssue(uint(id))
	if err != nil {
		respondErr(w, http.StatusNotFound, err.Error())
		return
	}
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(iss)
}

// 이슈 수정 핸들러
func (h *HTTPHandler) updateIssue(w http.ResponseWriter, r *http.Request) {
	// URL 파라미터에서 ID 추출
	idStr := chi.URLParam(r, "id")
	idI, err := strconv.Atoi(idStr)
	if err != nil {
		respondErr(w, http.StatusBadRequest, "invalid id")
		return
	}

	// 요청 바디 파싱
	var body struct {
		Title       *string            `json:"title"`
		Description *string            `json:"description"`
		Status      *model.IssueStatus `json:"status"`
		UserID      *uint              `json:"userId"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		respondErr(w, http.StatusBadRequest, "invalid json")
		return
	}

	// 서비스 호출
	iss, err := h.svc.UpdateIssue(uint(idI), body.Title, body.Description, body.Status, body.UserID)
	if err != nil {
		respondErr(w, http.StatusBadRequest, err.Error())
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(iss)
}
