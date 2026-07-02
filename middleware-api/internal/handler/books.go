package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/ehawk61/bookLister-ai/internal/model"
	"github.com/ehawk61/bookLister-ai/internal/repository"
	"github.com/ehawk61/bookLister-ai/internal/response"
)

type paginatedResponse struct {
	Data       []model.Book       `json:"data"`
	Pagination paginationMetadata `json:"pagination"`
}

type paginationMetadata struct {
	Page       int `json:"page"`
	PageSize   int `json:"page_size"`
	Total      int `json:"total"`
	TotalPages int `json:"total_pages"`
}

func parsePagination(r *http.Request) (page, pageSize int, err string) {
	page = 1
	pageSize = 20

	if v := r.URL.Query().Get("page"); v != "" {
		p, e := strconv.Atoi(v)
		if e != nil || p < 1 {
			return 0, 0, "page must be a positive integer"
		}
		page = p
	}

	if v := r.URL.Query().Get("page_size"); v != "" {
		ps, e := strconv.Atoi(v)
		if e != nil || ps < 1 || ps > 100 {
			return 0, 0, "page_size must be between 1 and 100"
		}
		pageSize = ps
	}

	return page, pageSize, ""
}

func writePaginatedResponse(w http.ResponseWriter, books []model.Book, page, pageSize, total int) {
	totalPages := 0
	if total > 0 {
		totalPages = (total + pageSize - 1) / pageSize
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(paginatedResponse{
		Data: books,
		Pagination: paginationMetadata{
			Page:       page,
			PageSize:   pageSize,
			Total:      total,
			TotalPages: totalPages,
		},
	})
}

func ListBooks(repo repository.BookRepository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		page, pageSize, errMsg := parsePagination(r)
		if errMsg != "" {
			response.WriteError(w, http.StatusBadRequest, response.ErrorResponse{
				ErrorCode:        "BAD_REQUEST",
				ErrorDescription: "Invalid pagination parameters",
				Details:          errMsg,
			})
			return
		}

		offset := (page - 1) * pageSize
		books, total := repo.GetAll(offset, pageSize)
		writePaginatedResponse(w, books, page, pageSize, total)
	}
}

func CreateBooks(repo repository.BookRepository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var raw json.RawMessage
		if err := json.NewDecoder(r.Body).Decode(&raw); err != nil {
			response.WriteError(w, http.StatusBadRequest, response.ErrorResponse{
				ErrorCode:        "BAD_REQUEST",
				ErrorDescription: "Invalid JSON body",
			})
			return
		}

		var books []model.Book
		isArray := len(raw) > 0 && raw[0] == '['

		if isArray {
			if err := json.Unmarshal(raw, &books); err != nil {
				response.WriteError(w, http.StatusBadRequest, response.ErrorResponse{
					ErrorCode:        "BAD_REQUEST",
					ErrorDescription: "Invalid JSON body",
					Details:          err.Error(),
				})
				return
			}
		} else {
			var b model.Book
			if err := json.Unmarshal(raw, &b); err != nil {
				response.WriteError(w, http.StatusBadRequest, response.ErrorResponse{
					ErrorCode:        "BAD_REQUEST",
					ErrorDescription: "Invalid JSON body",
					Details:          err.Error(),
				})
				return
			}
			books = []model.Book{b}
		}

		if len(books) == 0 {
			response.WriteError(w, http.StatusBadRequest, response.ErrorResponse{
				ErrorCode:        "BAD_REQUEST",
				ErrorDescription: "At least one book is required",
			})
			return
		}

		for i, b := range books {
			if details := b.Validate(); details != "" {
				response.WriteError(w, http.StatusUnprocessableEntity, response.ErrorResponse{
					ErrorCode:        "BOOK_VALIDATION_ERROR",
					ErrorDescription: "The book request is not valid",
					Details:          "book " + strconv.Itoa(i) + ": " + details,
				})
				return
			}
		}

		created := make([]model.Book, len(books))
		for i, b := range books {
			created[i] = repo.Create(b)
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		if isArray {
			json.NewEncoder(w).Encode(created)
		} else {
			json.NewEncoder(w).Encode(created[0])
		}
	}
}

func GetBook(repo repository.BookRepository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := r.PathValue("id")
		book, ok := repo.GetByID(id)
		if !ok {
			response.WriteError(w, http.StatusNotFound, response.ErrorResponse{
				ErrorCode:        "BOOK_NOT_FOUND",
				ErrorDescription: "Book does not exist",
			})
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(book)
	}
}

func UpdateBook(repo repository.BookRepository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := r.PathValue("id")

		var book model.Book
		if err := json.NewDecoder(r.Body).Decode(&book); err != nil {
			response.WriteError(w, http.StatusBadRequest, response.ErrorResponse{
				ErrorCode:        "BAD_REQUEST",
				ErrorDescription: "Invalid JSON body",
			})
			return
		}

		if details := book.Validate(); details != "" {
			response.WriteError(w, http.StatusUnprocessableEntity, response.ErrorResponse{
				ErrorCode:        "BOOK_VALIDATION_ERROR",
				ErrorDescription: "The book request is not valid",
				Details:          details,
			})
			return
		}

		if _, ok := repo.Update(id, book); !ok {
			response.WriteError(w, http.StatusNotFound, response.ErrorResponse{
				ErrorCode:        "BOOK_NOT_FOUND",
				ErrorDescription: "Book does not exist",
			})
			return
		}

		w.WriteHeader(http.StatusNoContent)
	}
}

func DeleteBook(repo repository.BookRepository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := r.PathValue("id")
		if !repo.Delete(id) {
			response.WriteError(w, http.StatusNotFound, response.ErrorResponse{
				ErrorCode:        "BOOK_NOT_FOUND",
				ErrorDescription: "Book does not exist",
			})
			return
		}

		w.WriteHeader(http.StatusNoContent)
	}
}

func SearchBooks(repo repository.BookRepository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		title := r.URL.Query().Get("title")
		author := r.URL.Query().Get("author")

		if title == "" && author == "" {
			response.WriteError(w, http.StatusBadRequest, response.ErrorResponse{
				ErrorCode:        "BAD_REQUEST",
				ErrorDescription: "At least one search parameter is required",
				Details:          "provide title and/or author query parameter",
			})
			return
		}

		page, pageSize, errMsg := parsePagination(r)
		if errMsg != "" {
			response.WriteError(w, http.StatusBadRequest, response.ErrorResponse{
				ErrorCode:        "BAD_REQUEST",
				ErrorDescription: "Invalid pagination parameters",
				Details:          errMsg,
			})
			return
		}

		offset := (page - 1) * pageSize
		books, total := repo.Search(title, author, offset, pageSize)
		writePaginatedResponse(w, books, page, pageSize, total)
	}
}
