package repository

import (
	"strings"
	"sync"

	"github.com/ehawk61/bookLister-ai/internal/model"
)

type MemoryBookRepository struct {
	mu    sync.RWMutex
	books map[string]model.Book
	order []string
}

func NewMemoryBookRepository() *MemoryBookRepository {
	return &MemoryBookRepository{
		books: make(map[string]model.Book),
	}
}

func (m *MemoryBookRepository) GetAll(offset, limit int) ([]model.Book, int) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	total := len(m.order)
	if offset >= total {
		return []model.Book{}, total
	}

	end := offset + limit
	if end > total {
		end = total
	}

	result := make([]model.Book, 0, end-offset)
	for _, id := range m.order[offset:end] {
		result = append(result, m.books[id])
	}
	return result, total
}

func (m *MemoryBookRepository) GetByID(id string) (model.Book, bool) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	book, ok := m.books[id]
	return book, ok
}

func (m *MemoryBookRepository) Create(book model.Book) model.Book {
	m.mu.Lock()
	defer m.mu.Unlock()

	book.ID = model.NewID()
	m.books[book.ID] = book
	m.order = append(m.order, book.ID)
	return book
}

func (m *MemoryBookRepository) Update(id string, book model.Book) (model.Book, bool) {
	m.mu.Lock()
	defer m.mu.Unlock()

	if _, ok := m.books[id]; !ok {
		return model.Book{}, false
	}

	book.ID = id
	m.books[id] = book
	return book, true
}

func (m *MemoryBookRepository) Delete(id string) bool {
	m.mu.Lock()
	defer m.mu.Unlock()

	if _, ok := m.books[id]; !ok {
		return false
	}

	delete(m.books, id)
	for i, oid := range m.order {
		if oid == id {
			m.order = append(m.order[:i], m.order[i+1:]...)
			break
		}
	}
	return true
}

func (m *MemoryBookRepository) Search(title, author string, offset, limit int) ([]model.Book, int) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	titleLower := strings.ToLower(title)
	authorLower := strings.ToLower(author)

	var matched []model.Book
	for _, id := range m.order {
		book := m.books[id]
		if title != "" && !strings.Contains(strings.ToLower(book.Title), titleLower) {
			continue
		}
		if author != "" && !strings.Contains(strings.ToLower(book.Author), authorLower) {
			continue
		}
		matched = append(matched, book)
	}

	total := len(matched)
	if offset >= total {
		return []model.Book{}, total
	}

	end := offset + limit
	if end > total {
		end = total
	}
	return matched[offset:end], total
}
