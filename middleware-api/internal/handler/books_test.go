package handler

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/ehawk61/bookLister-ai/internal/model"
	"github.com/ehawk61/bookLister-ai/internal/repository"
)

func seedRepo(repo *repository.MemoryBookRepository, books ...model.Book) []model.Book {
	created := make([]model.Book, len(books))
	for i, b := range books {
		created[i] = repo.Create(b)
	}
	return created
}

func TestListBooks(t *testing.T) {
	t.Run("empty repository", func(t *testing.T) {
		repo := repository.NewMemoryBookRepository()
		h := ListBooks(repo)

		req := httptest.NewRequest(http.MethodGet, "/api/books", nil)
		w := httptest.NewRecorder()
		h(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("status = %d, want %d", w.Code, http.StatusOK)
		}

		var resp paginatedResponse
		json.Unmarshal(w.Body.Bytes(), &resp)
		if len(resp.Data) != 0 {
			t.Errorf("len(data) = %d, want 0", len(resp.Data))
		}
		if resp.Pagination.Total != 0 {
			t.Errorf("total = %d, want 0", resp.Pagination.Total)
		}
	})

	t.Run("default pagination", func(t *testing.T) {
		repo := repository.NewMemoryBookRepository()
		seedRepo(repo,
			model.Book{Title: "A", Author: "X"},
			model.Book{Title: "B", Author: "Y"},
			model.Book{Title: "C", Author: "Z"},
		)
		h := ListBooks(repo)

		req := httptest.NewRequest(http.MethodGet, "/api/books", nil)
		w := httptest.NewRecorder()
		h(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("status = %d, want %d", w.Code, http.StatusOK)
		}

		var resp paginatedResponse
		json.Unmarshal(w.Body.Bytes(), &resp)
		if len(resp.Data) != 3 {
			t.Errorf("len(data) = %d, want 3", len(resp.Data))
		}
		if resp.Pagination.Total != 3 {
			t.Errorf("total = %d, want 3", resp.Pagination.Total)
		}
		if resp.Pagination.Page != 1 {
			t.Errorf("page = %d, want 1", resp.Pagination.Page)
		}
	})

	t.Run("custom page and page_size", func(t *testing.T) {
		repo := repository.NewMemoryBookRepository()
		for i := 0; i < 25; i++ {
			repo.Create(model.Book{Title: "Book", Author: "Author"})
		}
		h := ListBooks(repo)

		req := httptest.NewRequest(http.MethodGet, "/api/books?page=2&page_size=10", nil)
		w := httptest.NewRecorder()
		h(w, req)

		var resp paginatedResponse
		json.Unmarshal(w.Body.Bytes(), &resp)
		if len(resp.Data) != 10 {
			t.Errorf("len(data) = %d, want 10", len(resp.Data))
		}
		if resp.Pagination.Total != 25 {
			t.Errorf("total = %d, want 25", resp.Pagination.Total)
		}
		if resp.Pagination.TotalPages != 3 {
			t.Errorf("total_pages = %d, want 3", resp.Pagination.TotalPages)
		}
	})

	t.Run("invalid page", func(t *testing.T) {
		repo := repository.NewMemoryBookRepository()
		h := ListBooks(repo)

		req := httptest.NewRequest(http.MethodGet, "/api/books?page=-1", nil)
		w := httptest.NewRecorder()
		h(w, req)

		if w.Code != http.StatusBadRequest {
			t.Errorf("status = %d, want %d", w.Code, http.StatusBadRequest)
		}
	})

	t.Run("invalid page_size", func(t *testing.T) {
		repo := repository.NewMemoryBookRepository()
		h := ListBooks(repo)

		req := httptest.NewRequest(http.MethodGet, "/api/books?page_size=999", nil)
		w := httptest.NewRecorder()
		h(w, req)

		if w.Code != http.StatusBadRequest {
			t.Errorf("status = %d, want %d", w.Code, http.StatusBadRequest)
		}
	})
}

func TestCreateBooks(t *testing.T) {
	t.Run("single valid book", func(t *testing.T) {
		repo := repository.NewMemoryBookRepository()
		h := CreateBooks(repo)

		body := `{"title":"Dune","author":"Frank Herbert"}`
		req := httptest.NewRequest(http.MethodPost, "/api/books", strings.NewReader(body))
		w := httptest.NewRecorder()
		h(w, req)

		if w.Code != http.StatusCreated {
			t.Errorf("status = %d, want %d", w.Code, http.StatusCreated)
		}

		var book model.Book
		json.Unmarshal(w.Body.Bytes(), &book)
		if book.ID == "" {
			t.Error("expected generated ID")
		}
		if book.Title != "Dune" {
			t.Errorf("title = %q, want %q", book.Title, "Dune")
		}

		stored, ok := repo.GetByID(book.ID)
		if !ok {
			t.Fatal("book not found in repository after creation")
		}
		if stored.Title != "Dune" {
			t.Errorf("stored title = %q, want %q", stored.Title, "Dune")
		}
	})

	t.Run("array of books", func(t *testing.T) {
		repo := repository.NewMemoryBookRepository()
		h := CreateBooks(repo)

		body := `[{"title":"Dune","author":"Frank Herbert"},{"title":"Foundation","author":"Isaac Asimov"}]`
		req := httptest.NewRequest(http.MethodPost, "/api/books", strings.NewReader(body))
		w := httptest.NewRecorder()
		h(w, req)

		if w.Code != http.StatusCreated {
			t.Errorf("status = %d, want %d", w.Code, http.StatusCreated)
		}

		var books []model.Book
		json.Unmarshal(w.Body.Bytes(), &books)
		if len(books) != 2 {
			t.Fatalf("len(books) = %d, want 2", len(books))
		}
		if books[0].ID == "" || books[1].ID == "" {
			t.Error("expected generated IDs for all books")
		}
	})

	t.Run("missing title", func(t *testing.T) {
		repo := repository.NewMemoryBookRepository()
		h := CreateBooks(repo)

		body := `{"author":"Frank Herbert"}`
		req := httptest.NewRequest(http.MethodPost, "/api/books", strings.NewReader(body))
		w := httptest.NewRecorder()
		h(w, req)

		if w.Code != http.StatusUnprocessableEntity {
			t.Errorf("status = %d, want %d", w.Code, http.StatusUnprocessableEntity)
		}

		var errResp map[string]string
		json.Unmarshal(w.Body.Bytes(), &errResp)
		if errResp["error_code"] != "BOOK_VALIDATION_ERROR" {
			t.Errorf("error_code = %q, want BOOK_VALIDATION_ERROR", errResp["error_code"])
		}
	})

	t.Run("missing author", func(t *testing.T) {
		repo := repository.NewMemoryBookRepository()
		h := CreateBooks(repo)

		body := `{"title":"Dune"}`
		req := httptest.NewRequest(http.MethodPost, "/api/books", strings.NewReader(body))
		w := httptest.NewRecorder()
		h(w, req)

		if w.Code != http.StatusUnprocessableEntity {
			t.Errorf("status = %d, want %d", w.Code, http.StatusUnprocessableEntity)
		}
	})

	t.Run("malformed JSON", func(t *testing.T) {
		repo := repository.NewMemoryBookRepository()
		h := CreateBooks(repo)

		req := httptest.NewRequest(http.MethodPost, "/api/books", strings.NewReader("{invalid"))
		w := httptest.NewRecorder()
		h(w, req)

		if w.Code != http.StatusBadRequest {
			t.Errorf("status = %d, want %d", w.Code, http.StatusBadRequest)
		}
	})

	t.Run("empty array", func(t *testing.T) {
		repo := repository.NewMemoryBookRepository()
		h := CreateBooks(repo)

		req := httptest.NewRequest(http.MethodPost, "/api/books", strings.NewReader("[]"))
		w := httptest.NewRecorder()
		h(w, req)

		if w.Code != http.StatusBadRequest {
			t.Errorf("status = %d, want %d", w.Code, http.StatusBadRequest)
		}
	})
}

func TestGetBook(t *testing.T) {
	t.Run("found", func(t *testing.T) {
		repo := repository.NewMemoryBookRepository()
		created := seedRepo(repo, model.Book{Title: "Dune", Author: "Frank Herbert"})
		h := GetBook(repo)

		req := httptest.NewRequest(http.MethodGet, "/api/books/"+created[0].ID, nil)
		req.SetPathValue("id", created[0].ID)
		w := httptest.NewRecorder()
		h(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("status = %d, want %d", w.Code, http.StatusOK)
		}

		var book model.Book
		json.Unmarshal(w.Body.Bytes(), &book)
		if book.Title != "Dune" {
			t.Errorf("title = %q, want %q", book.Title, "Dune")
		}
	})

	t.Run("not found", func(t *testing.T) {
		repo := repository.NewMemoryBookRepository()
		h := GetBook(repo)

		req := httptest.NewRequest(http.MethodGet, "/api/books/nonexistent", nil)
		req.SetPathValue("id", "nonexistent")
		w := httptest.NewRecorder()
		h(w, req)

		if w.Code != http.StatusNotFound {
			t.Errorf("status = %d, want %d", w.Code, http.StatusNotFound)
		}

		var errResp map[string]string
		json.Unmarshal(w.Body.Bytes(), &errResp)
		if errResp["error_code"] != "BOOK_NOT_FOUND" {
			t.Errorf("error_code = %q, want BOOK_NOT_FOUND", errResp["error_code"])
		}
	})
}

func TestUpdateBook(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		repo := repository.NewMemoryBookRepository()
		created := seedRepo(repo, model.Book{Title: "Dune", Author: "Frank Herbert"})
		h := UpdateBook(repo)

		body := `{"title":"Dune Messiah","author":"Frank Herbert"}`
		req := httptest.NewRequest(http.MethodPut, "/api/books/"+created[0].ID, strings.NewReader(body))
		req.SetPathValue("id", created[0].ID)
		w := httptest.NewRecorder()
		h(w, req)

		if w.Code != http.StatusNoContent {
			t.Errorf("status = %d, want %d", w.Code, http.StatusNoContent)
		}

		if w.Body.Len() != 0 {
			t.Errorf("body = %q, want empty", w.Body.String())
		}

		updated, _ := repo.GetByID(created[0].ID)
		if updated.Title != "Dune Messiah" {
			t.Errorf("title = %q, want %q", updated.Title, "Dune Messiah")
		}
	})

	t.Run("not found", func(t *testing.T) {
		repo := repository.NewMemoryBookRepository()
		h := UpdateBook(repo)

		body := `{"title":"X","author":"Y"}`
		req := httptest.NewRequest(http.MethodPut, "/api/books/nonexistent", strings.NewReader(body))
		req.SetPathValue("id", "nonexistent")
		w := httptest.NewRecorder()
		h(w, req)

		if w.Code != http.StatusNotFound {
			t.Errorf("status = %d, want %d", w.Code, http.StatusNotFound)
		}
	})

	t.Run("validation error", func(t *testing.T) {
		repo := repository.NewMemoryBookRepository()
		created := seedRepo(repo, model.Book{Title: "Dune", Author: "Frank Herbert"})
		h := UpdateBook(repo)

		body := `{"title":"","author":"Frank Herbert"}`
		req := httptest.NewRequest(http.MethodPut, "/api/books/"+created[0].ID, strings.NewReader(body))
		req.SetPathValue("id", created[0].ID)
		w := httptest.NewRecorder()
		h(w, req)

		if w.Code != http.StatusUnprocessableEntity {
			t.Errorf("status = %d, want %d", w.Code, http.StatusUnprocessableEntity)
		}

		var errResp map[string]string
		json.Unmarshal(w.Body.Bytes(), &errResp)
		if errResp["error_code"] != "BOOK_VALIDATION_ERROR" {
			t.Errorf("error_code = %q, want BOOK_VALIDATION_ERROR", errResp["error_code"])
		}
	})

	t.Run("malformed JSON", func(t *testing.T) {
		repo := repository.NewMemoryBookRepository()
		created := seedRepo(repo, model.Book{Title: "Dune", Author: "Frank Herbert"})
		h := UpdateBook(repo)

		req := httptest.NewRequest(http.MethodPut, "/api/books/"+created[0].ID, strings.NewReader("{bad"))
		req.SetPathValue("id", created[0].ID)
		w := httptest.NewRecorder()
		h(w, req)

		if w.Code != http.StatusBadRequest {
			t.Errorf("status = %d, want %d", w.Code, http.StatusBadRequest)
		}
	})
}

func TestDeleteBook(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		repo := repository.NewMemoryBookRepository()
		created := seedRepo(repo, model.Book{Title: "Dune", Author: "Frank Herbert"})
		h := DeleteBook(repo)

		req := httptest.NewRequest(http.MethodDelete, "/api/books/"+created[0].ID, nil)
		req.SetPathValue("id", created[0].ID)
		w := httptest.NewRecorder()
		h(w, req)

		if w.Code != http.StatusNoContent {
			t.Errorf("status = %d, want %d", w.Code, http.StatusNoContent)
		}

		_, ok := repo.GetByID(created[0].ID)
		if ok {
			t.Error("expected book to be deleted")
		}
	})

	t.Run("not found", func(t *testing.T) {
		repo := repository.NewMemoryBookRepository()
		h := DeleteBook(repo)

		req := httptest.NewRequest(http.MethodDelete, "/api/books/nonexistent", nil)
		req.SetPathValue("id", "nonexistent")
		w := httptest.NewRecorder()
		h(w, req)

		if w.Code != http.StatusNotFound {
			t.Errorf("status = %d, want %d", w.Code, http.StatusNotFound)
		}
	})
}

func TestSearchBooks(t *testing.T) {
	setupRepo := func() *repository.MemoryBookRepository {
		repo := repository.NewMemoryBookRepository()
		seedRepo(repo,
			model.Book{Title: "Dune", Author: "Frank Herbert"},
			model.Book{Title: "Neuromancer", Author: "William Gibson"},
			model.Book{Title: "Foundation", Author: "Isaac Asimov"},
		)
		return repo
	}

	t.Run("by title", func(t *testing.T) {
		repo := setupRepo()
		h := SearchBooks(repo)

		req := httptest.NewRequest(http.MethodGet, "/api/books/q?title=Dune", nil)
		w := httptest.NewRecorder()
		h(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("status = %d, want %d", w.Code, http.StatusOK)
		}

		var resp paginatedResponse
		json.Unmarshal(w.Body.Bytes(), &resp)
		if len(resp.Data) != 1 {
			t.Errorf("len(data) = %d, want 1", len(resp.Data))
		}
		if resp.Data[0].Title != "Dune" {
			t.Errorf("title = %q, want %q", resp.Data[0].Title, "Dune")
		}
	})

	t.Run("by author", func(t *testing.T) {
		repo := setupRepo()
		h := SearchBooks(repo)

		req := httptest.NewRequest(http.MethodGet, "/api/books/q?author=Gibson", nil)
		w := httptest.NewRecorder()
		h(w, req)

		var resp paginatedResponse
		json.Unmarshal(w.Body.Bytes(), &resp)
		if len(resp.Data) != 1 {
			t.Errorf("len(data) = %d, want 1", len(resp.Data))
		}
	})

	t.Run("by title and author", func(t *testing.T) {
		repo := setupRepo()
		h := SearchBooks(repo)

		req := httptest.NewRequest(http.MethodGet, "/api/books/q?title=Dune&author=Herbert", nil)
		w := httptest.NewRecorder()
		h(w, req)

		var resp paginatedResponse
		json.Unmarshal(w.Body.Bytes(), &resp)
		if len(resp.Data) != 1 {
			t.Errorf("len(data) = %d, want 1", len(resp.Data))
		}
	})

	t.Run("no matches", func(t *testing.T) {
		repo := setupRepo()
		h := SearchBooks(repo)

		req := httptest.NewRequest(http.MethodGet, "/api/books/q?title=Nonexistent", nil)
		w := httptest.NewRecorder()
		h(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("status = %d, want %d", w.Code, http.StatusOK)
		}

		var resp paginatedResponse
		json.Unmarshal(w.Body.Bytes(), &resp)
		if len(resp.Data) != 0 {
			t.Errorf("len(data) = %d, want 0", len(resp.Data))
		}
	})

	t.Run("no query params", func(t *testing.T) {
		repo := setupRepo()
		h := SearchBooks(repo)

		req := httptest.NewRequest(http.MethodGet, "/api/books/q", nil)
		w := httptest.NewRecorder()
		h(w, req)

		if w.Code != http.StatusBadRequest {
			t.Errorf("status = %d, want %d", w.Code, http.StatusBadRequest)
		}
	})

	t.Run("with pagination", func(t *testing.T) {
		repo := repository.NewMemoryBookRepository()
		for i := 0; i < 15; i++ {
			repo.Create(model.Book{Title: "Go Programming", Author: "Author"})
		}
		h := SearchBooks(repo)

		req := httptest.NewRequest(http.MethodGet, "/api/books/q?title=Go&page=2&page_size=10", nil)
		w := httptest.NewRecorder()
		h(w, req)

		var resp paginatedResponse
		json.Unmarshal(w.Body.Bytes(), &resp)
		if len(resp.Data) != 5 {
			t.Errorf("len(data) = %d, want 5", len(resp.Data))
		}
		if resp.Pagination.Total != 15 {
			t.Errorf("total = %d, want 15", resp.Pagination.Total)
		}
		if resp.Pagination.TotalPages != 2 {
			t.Errorf("total_pages = %d, want 2", resp.Pagination.TotalPages)
		}
	})
}
