package repository

import "github.com/ehawk61/bookLister-ai/internal/model"

type BookRepository interface {
	GetAll(offset, limit int) ([]model.Book, int)
	GetByID(id string) (model.Book, bool)
	Create(book model.Book) model.Book
	Update(id string, book model.Book) (model.Book, bool)
	Delete(id string) bool
	Search(title, author string, offset, limit int) ([]model.Book, int)
}
