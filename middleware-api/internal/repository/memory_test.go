package repository

import (
	"testing"

	"github.com/ehawk61/bookLister-ai/internal/model"
)

func TestMemoryBookRepository_Create(t *testing.T) {
	repo := NewMemoryBookRepository()

	book := repo.Create(model.Book{Title: "Dune", Author: "Frank Herbert"})

	if book.ID == "" {
		t.Fatal("expected generated ID, got empty string")
	}
	if book.Title != "Dune" {
		t.Errorf("title = %q, want %q", book.Title, "Dune")
	}
	if book.Author != "Frank Herbert" {
		t.Errorf("author = %q, want %q", book.Author, "Frank Herbert")
	}
}

func TestMemoryBookRepository_GetByID(t *testing.T) {
	repo := NewMemoryBookRepository()
	created := repo.Create(model.Book{Title: "Dune", Author: "Frank Herbert"})

	t.Run("found", func(t *testing.T) {
		book, ok := repo.GetByID(created.ID)
		if !ok {
			t.Fatal("expected book to be found")
		}
		if book.Title != "Dune" {
			t.Errorf("title = %q, want %q", book.Title, "Dune")
		}
	})

	t.Run("not found", func(t *testing.T) {
		_, ok := repo.GetByID("nonexistent")
		if ok {
			t.Fatal("expected book not to be found")
		}
	})
}

func TestMemoryBookRepository_GetAll(t *testing.T) {
	repo := NewMemoryBookRepository()

	t.Run("empty repository", func(t *testing.T) {
		books, total := repo.GetAll(0, 10)
		if total != 0 {
			t.Errorf("total = %d, want 0", total)
		}
		if len(books) != 0 {
			t.Errorf("len(books) = %d, want 0", len(books))
		}
	})

	repo.Create(model.Book{Title: "Book A", Author: "Author A"})
	repo.Create(model.Book{Title: "Book B", Author: "Author B"})
	repo.Create(model.Book{Title: "Book C", Author: "Author C"})

	t.Run("all books", func(t *testing.T) {
		books, total := repo.GetAll(0, 10)
		if total != 3 {
			t.Errorf("total = %d, want 3", total)
		}
		if len(books) != 3 {
			t.Errorf("len(books) = %d, want 3", len(books))
		}
	})

	t.Run("pagination", func(t *testing.T) {
		books, total := repo.GetAll(1, 1)
		if total != 3 {
			t.Errorf("total = %d, want 3", total)
		}
		if len(books) != 1 {
			t.Errorf("len(books) = %d, want 1", len(books))
		}
		if books[0].Title != "Book B" {
			t.Errorf("title = %q, want %q", books[0].Title, "Book B")
		}
	})

	t.Run("offset beyond total", func(t *testing.T) {
		books, total := repo.GetAll(10, 5)
		if total != 3 {
			t.Errorf("total = %d, want 3", total)
		}
		if len(books) != 0 {
			t.Errorf("len(books) = %d, want 0", len(books))
		}
	})
}

func TestMemoryBookRepository_Update(t *testing.T) {
	repo := NewMemoryBookRepository()
	created := repo.Create(model.Book{Title: "Dune", Author: "Frank Herbert"})

	t.Run("existing book", func(t *testing.T) {
		updated, ok := repo.Update(created.ID, model.Book{Title: "Dune Messiah", Author: "Frank Herbert"})
		if !ok {
			t.Fatal("expected update to succeed")
		}
		if updated.Title != "Dune Messiah" {
			t.Errorf("title = %q, want %q", updated.Title, "Dune Messiah")
		}
		if updated.ID != created.ID {
			t.Errorf("ID = %q, want %q", updated.ID, created.ID)
		}

		fetched, _ := repo.GetByID(created.ID)
		if fetched.Title != "Dune Messiah" {
			t.Errorf("fetched title = %q, want %q", fetched.Title, "Dune Messiah")
		}
	})

	t.Run("nonexistent book", func(t *testing.T) {
		_, ok := repo.Update("nonexistent", model.Book{Title: "X", Author: "Y"})
		if ok {
			t.Fatal("expected update to fail for nonexistent book")
		}
	})
}

func TestMemoryBookRepository_Delete(t *testing.T) {
	repo := NewMemoryBookRepository()
	created := repo.Create(model.Book{Title: "Dune", Author: "Frank Herbert"})

	t.Run("existing book", func(t *testing.T) {
		if !repo.Delete(created.ID) {
			t.Fatal("expected delete to succeed")
		}
		_, ok := repo.GetByID(created.ID)
		if ok {
			t.Fatal("expected book to be gone after delete")
		}
		books, total := repo.GetAll(0, 10)
		if total != 0 {
			t.Errorf("total = %d, want 0", total)
		}
		if len(books) != 0 {
			t.Errorf("len(books) = %d, want 0", len(books))
		}
	})

	t.Run("nonexistent book", func(t *testing.T) {
		if repo.Delete("nonexistent") {
			t.Fatal("expected delete to fail for nonexistent book")
		}
	})
}

func TestMemoryBookRepository_Search(t *testing.T) {
	repo := NewMemoryBookRepository()
	repo.Create(model.Book{Title: "Dune", Author: "Frank Herbert"})
	repo.Create(model.Book{Title: "Neuromancer", Author: "William Gibson"})
	repo.Create(model.Book{Title: "Foundation", Author: "Isaac Asimov"})

	tests := []struct {
		name      string
		title     string
		author    string
		wantCount int
	}{
		{"by title exact", "Dune", "", 1},
		{"by title substring", "un", "", 2},
		{"by title case insensitive", "dune", "", 1},
		{"by author", "", "Gibson", 1},
		{"by title and author", "Dune", "Herbert", 1},
		{"by title and author no match", "Dune", "Gibson", 0},
		{"no matches", "Xyz", "", 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			books, total := repo.Search(tt.title, tt.author, 0, 100)
			if total != tt.wantCount {
				t.Errorf("total = %d, want %d", total, tt.wantCount)
			}
			if len(books) != tt.wantCount {
				t.Errorf("len(books) = %d, want %d", len(books), tt.wantCount)
			}
		})
	}

	t.Run("search pagination", func(t *testing.T) {
		books, total := repo.Search("un", "", 0, 1)
		if total != 2 {
			t.Errorf("total = %d, want 2", total)
		}
		if len(books) != 1 {
			t.Errorf("len(books) = %d, want 1", len(books))
		}
	})
}
