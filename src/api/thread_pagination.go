package api

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"
)

// Thread represents a discussion thread.
type Thread struct {
	ID        int64     `json:"id"`
	Title     string    `json:"title"`
	Content   string    `json:"content"` // markdown content
	CreatedAt time.Time `json:"created_at"`
}

// pageSize is the number of threads returned per page.
const pageSize = 20

// mockThreads is an in‑memory data source used for this example.
// In a real implementation this would be replaced by DB calls.
var mockThreads []Thread

func init() {
	// Populate mockThreads with dummy data for demonstration.
	for i := 1; i <= 105; i++ {
		mockThreads = append(mockThreads, Thread{
			ID:        int64(i),
			Title:     "Thread #" + strconv.Itoa(i),
			Content:   "# Header\n\nSample **markdown** content for thread " + strconv.Itoa(i),
			CreatedAt: time.Now().Add(-time.Duration(i) * time.Hour),
		})
	}
}

// PaginatedResponse is the JSON payload returned by the pagination endpoint.
type PaginatedResponse struct {
	Page       int      `json:"page"`
	TotalPages int      `json:"total_pages"`
	NextPage   *int     `json:"next_page,omitempty"`
	PrevPage   *int     `json:"prev_page,omitempty"`
	Threads    []Thread `json:"threads"`
}

// ThreadPaginationHandler handles GET /threads requests with pagination.
// Query parameters:
//   - page (optional, default 1): the page number to retrieve.
func ThreadPaginationHandler(w http.ResponseWriter, r *http.Request) {
	// Parse page query param.
	page := 1
	if pStr := r.URL.Query().Get("page"); pStr != "" {
		if p, err := strconv.Atoi(pStr); err == nil && p > 0 {
			page = p
		}
	}

	total := len(mockThreads)
	totalPages := (total + pageSize - 1) / pageSize // ceiling division

	// Clamp page within valid range.
	if page > totalPages {
		page = totalPages
	}
	if page < 1 {
		page = 1
	}

	start := (page - 1) * pageSize
	end := start + pageSize
	if end > total {
		end = total
	}
	pageThreads := mockThreads[start:end]

	// Build navigation pointers.
	var nextPage *int
	var prevPage *int
	if page < totalPages {
		n := page + 1
		nextPage = &n
	}
	if page > 1 {
		p := page - 1
		prevPage = &p
	}

	resp := PaginatedResponse{
		Page:       page,
		TotalPages: totalPages,
		NextPage:   nextPage,
		PrevPage:   prevPage,
		Threads:    pageThreads,
	}

	w.Header().Set("Content-Type", "application/json")
	enc := json.NewEncoder(w)
	if err := enc.Encode(resp); err != nil {
		http.Error(w, "failed to encode response", http.StatusInternalServerError)
	}
}

// RegisterThreadRoutes registers the thread pagination endpoint onto the provided mux.
func RegisterThreadRoutes(mux *http.ServeMux) {
	mux.HandleFunc("/threads", ThreadPaginationHandler)
}