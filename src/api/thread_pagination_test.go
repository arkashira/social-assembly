package api

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"
)

func TestThreadPaginationHandler_FirstPage(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/threads?page=1", nil)
	rr := httptest.NewRecorder()

	ThreadPaginationHandler(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, status)
	}

	var resp PaginatedResponse
	if err := json.NewDecoder(rr.Body).Decode(&resp); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if resp.Page != 1 {
		t.Errorf("expected page 1, got %d", resp.Page)
	}
	if resp.TotalPages == 0 {
		t.Errorf("expected total pages > 0")
	}
	if len(resp.Threads) != pageSize {
		t.Errorf("expected %d threads, got %d", pageSize, len(resp.Threads))
	}
	if resp.NextPage == nil || *resp.NextPage != 2 {
		t.Errorf("expected next_page 2, got %+v", resp.NextPage)
	}
	if resp.PrevPage != nil {
		t.Errorf("expected prev_page nil on first page, got %+v", resp.PrevPage)
	}
}

func TestThreadPaginationHandler_LastPage(t *testing.T) {
	// Compute last page based on mock data.
	total := len(mockThreads)
	lastPage := (total + pageSize - 1) / pageSize

	req := httptest.NewRequest(http.MethodGet, "/threads?page="+strconv.Itoa(lastPage), nil)
	rr := httptest.NewRecorder()

	ThreadPaginationHandler(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, status)
	}

	var resp PaginatedResponse
	if err := json.NewDecoder(rr.Body).Decode(&resp); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if resp.Page != lastPage {
		t.Errorf("expected page %d, got %d", lastPage, resp.Page)
	}
	if resp.NextPage != nil {
		t.Errorf("expected next_page nil on last page, got %+v", resp.NextPage)
	}
	if resp.PrevPage == nil || *resp.PrevPage != lastPage-1 {
		t.Errorf("expected prev_page %d, got %+v", lastPage-1, resp.PrevPage)
	}
	// The last page may have fewer than pageSize items.
	expectedCount := total - (pageSize * (lastPage - 1))
	if len(resp.Threads) != expectedCount {
		t.Errorf("expected %d threads on last page, got %d", expectedCount, len(resp.Threads))
	}
}