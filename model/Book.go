package model

import "time"

// Book represents a book from a bookstore system.
type Book struct {
	ID           int64 `bun:"id,pk,autoincrement" gorm:"primary_key"`
	ISBN         string
	Title        string
	Author       string
	Genre        string
	Quantity     int
	PublicizedAt time.Time
}

func NewBooks(quantity int) []*Book {
	books := make([]*Book, quantity)
	for i := 0; i < quantity; i++ {
		books[i] = NewBook()
	}
	return books
}

func NewBooksNoPtr(quantity int) []Book {
	books := make([]Book, quantity)
	for i := 0; i < quantity; i++ {
		books[i] = NewBookNoPtr()
	}
	return books
}

func NewBookNoPtr() Book {
	return Book{
		ISBN:         "978-3-16-148410-1",
		Title:        "Learning Go: An Idiomatic Approach to Real-World Go Programming",
		Author:       "Jon Bodner",
		Genre:        "Programming",
		Quantity:     20,
		PublicizedAt: time.Date(2022, time.January, 1, 0, 0, 0, 0, time.UTC),
	}
}

func NewBook() *Book {
	return &Book{
		ISBN:         "978-3-16-148410-1",
		Title:        "Learning Go: An Idiomatic Approach to Real-World Go Programming",
		Author:       "Jon Bodner",
		Genre:        "Programming",
		Quantity:     20,
		PublicizedAt: time.Date(2022, time.January, 1, 0, 0, 0, 0, time.UTC),
	}
}

func Chunk(input []*Book, batchSize int) [][]*Book {
	var result [][]*Book
	for i := 0; i < len(input); i += batchSize {
		end := i + batchSize
		if end > len(input) {
			end = len(input)
		}
		result = append(result, input[i:end])
	}
	return result
}
