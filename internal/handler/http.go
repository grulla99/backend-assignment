package handler

import (
    "encoding/json"
    "net/http"
    "strconv"

    "github.com/go-chi/chi/v5"

    "backed-assignment-aro/internal/model"
    "backed-assignment-aro/internal/service"
)

// HTTPHandler wires routes.
type HTTPHandler struct {
    svc *service.IssueService
}

func NewHTTPHandler(s *service.IssueService) *HTTPHandler {
    return &HTTPHandler{svc: s}
}

func (h *HTTPHandler) Register(r chi.Router) {
    r.Post("/issue", h.createIssue)
    r.Get("/issues", h.listIssues)
    r.Get("/issue/{id}", h.getIssue)
    r.Patch("/issue/{id}", h.updateIssue)
}

func respondErr(w http.ResponseWriter, code int, msg string) {
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(code)
    _ = json.NewEncoder(w).Encode(map[string]interface{}{
        "error": msg,
        "code":  code,
    })
}

func (h *HTTPHandler) createIssue(w http.ResponseWriter, r *http.Request) {
    var body struct {
        Title       string `json:"title"`
        Description string `json:"description"`
        UserID      *uint  `json:"userId"`
    }
    if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
        respondErr(w, http.StatusBadRequest, "invalid json")
        return
    }
    iss, err := h.svc.CreateIssue(body.Title, body.Description, body.UserID)
    if err != nil {
        respondErr(w, http.StatusBadRequest, err.Error())
        return
    }
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusCreated)
    _ = json.NewEncoder(w).Encode(iss)
}

func (h *HTTPHandler) listIssues(w http.ResponseWriter, r *http.Request) {
    q := r.URL.Query().Get("status")
    var status *model.IssueStatus
    if q != "" {
        s := model.IssueStatus(q)
        status = &s
    }
    list := h.svc.ListIssues(status)
    resp := map[string]interface{}{"issues": list}
    w.Header().Set("Content-Type", "application/json")
    _ = json.NewEncoder(w).Encode(resp)
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

func (h *HTTPHandler) updateIssue(w http.ResponseWriter, r *http.Request) {
    idStr := chi.URLParam(r, "id")
    idI, err := strconv.Atoi(idStr)
    if err != nil {
        respondErr(w, http.StatusBadRequest, "invalid id")
        return
    }
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
    iss, err := h.svc.UpdateIssue(uint(idI), body.Title, body.Description, body.Status, body.UserID)
    if err != nil {
        respondErr(w, http.StatusBadRequest, err.Error())
        return
    }
    w.Header().Set("Content-Type", "application/json")
    _ = json.NewEncoder(w).Encode(iss)
}
