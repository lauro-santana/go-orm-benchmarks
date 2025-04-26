package benchmark

import (
	"database/sql"
	"fmt"
	"strings"
	"testing"

	"github.com/lauro-santana/golang-orm-benchmarks/benchmark/utils"
	"github.com/lauro-santana/golang-orm-benchmarks/model"
)

type RawBenchmark struct {
	db *sql.DB
}

func NewRawBenchmark() Benchmark {
	return &RawBenchmark{}
}

func (r *RawBenchmark) Init() error {
	var err error
	r.db, err = sql.Open("pgx", utils.PostgresDSN)
	return err
}

func (r *RawBenchmark) Close() error {
	return r.db.Close()
}

func (r *RawBenchmark) Insert(b *testing.B) {
	book := model.NewBook()

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_, err := r.db.Exec(utils.InsertQuery,
			book.ISBN, book.Title, book.Author, book.Genre, book.Quantity, book.PublicizedAt)

		b.StopTimer()
		if err != nil {
			b.Error(err)
		}
		b.StartTimer()
	}
}

func (r *RawBenchmark) InsertBulk(b *testing.B) {
	books := model.NewBooks(utils.BulkInsertNumber)

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		err := r.doInsertBulk(books)

		if err != nil {
			b.Error(err)
		}
	}
}

func (r *RawBenchmark) Update(b *testing.B) {
	book := model.NewBook()
	var id int64
	err := r.db.QueryRow(utils.InsertReturningIDQuery,
		book.ISBN, book.Title, book.Author, book.Genre, book.Quantity, book.PublicizedAt).Scan(&id)
	if err != nil {
		b.Error(err)
	}

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_, err = r.db.Exec(utils.UpdateQuery,
			book.ISBN, book.Title, book.Author, book.Genre, book.Quantity, book.PublicizedAt, id)

		if err != nil {
			b.Error(err)
		}
	}
}

func (r *RawBenchmark) Delete(b *testing.B) {
	n := b.N
	book := model.NewBook()
	bookIDs := make([]int64, 0, n)
	for i := 0; i < n; i++ {
		var id int64
		err := r.db.QueryRow(utils.InsertReturningIDQuery,
			book.ISBN, book.Title, book.Author, book.Genre, book.Quantity, book.PublicizedAt).Scan(&id)
		if err != nil {
			b.Error(err)
		}
		bookIDs = append(bookIDs, id)
	}

	b.ReportAllocs()
	b.ResetTimer()

	var bookID int64
	for i := 0; i < b.N; i++ {
		b.StopTimer()
		bookID = bookIDs[i]
		b.StartTimer()

		_, err := r.db.Exec(utils.DeleteQuery, bookID)

		if err != nil {
			b.Error(err)
		}
	}
}

func (r *RawBenchmark) FindByID(b *testing.B) {
	book := model.NewBook()
	var id int64
	err := r.db.QueryRow(utils.InsertReturningIDQuery,
		book.ISBN, book.Title, book.Author, book.Genre, book.Quantity, book.PublicizedAt).Scan(&id)
	if err != nil {
		b.Error(err)
	}

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		for range utils.FindOneLoop {
			var foundBook model.Book
			err := r.db.QueryRow(utils.SelectByIDQuery, id).Scan(
				&foundBook.ID,
				&foundBook.ISBN,
				&foundBook.Title,
				&foundBook.Author,
				&foundBook.Genre,
				&foundBook.Quantity,
				&foundBook.PublicizedAt,
			)

			// checking the error will count on raw benchmarks
			if err != nil {
				b.Error(err)
			}
		}
	}
}

var booksPage []model.Book

func (r *RawBenchmark) FindPage(b *testing.B) {
	books := model.NewBooks(utils.BulkInsertPageNumber)
	batches := model.Chunk(books, utils.BatchSize)
	for _, batch := range batches {
		if err := r.doInsertBulk(batch); err != nil {
			b.Error(err)
		}
	}

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		for s := 0; s < utils.BulkInsertPageNumber; s = s + utils.PageSize {
			// making slices will count on raw benchmarks
			booksPage = make([]model.Book, 0, utils.PageSize)

			rows, err := r.db.Query(utils.SelectPaginatingQuery, s, utils.PageSize)

			// checking the error will count on raw benchmarks
			if err != nil {
				b.Error(err)
			}

			for rows.Next() {
				var book model.Book
				if err = rows.Scan(
					&book.ID,
					&book.ISBN,
					&book.Title,
					&book.Author,
					&book.Genre,
					&book.Quantity,
					&book.PublicizedAt,
				); err != nil {
					b.Error(err)
				}
				booksPage = append(booksPage, book)
			}
		}
	}
}

func (r *RawBenchmark) doInsertBulk(books []*model.Book) error {
	valueStrings := make([]string, 0, len(books))
	valueArgs := make([]interface{}, 0, len(books)*6)

	start := 1

	for _, book := range books {
		placeholders := make([]string, 0, 6)
		for i := 0; i < 6; i++ {
			placeholders = append(placeholders, fmt.Sprintf("$%d", start))
			start++
		}
		valueStrings = append(valueStrings, "("+strings.Join(placeholders, ",")+")")
		valueArgs = append(valueArgs, book.ISBN)
		valueArgs = append(valueArgs, book.Title)
		valueArgs = append(valueArgs, book.Author)
		valueArgs = append(valueArgs, book.Genre)
		valueArgs = append(valueArgs, book.Quantity)
		valueArgs = append(valueArgs, book.PublicizedAt)
	}

	query := fmt.Sprintf(utils.InsertBulkQuery, strings.Join(valueStrings, ","))

	_, err := r.db.Exec(query, valueArgs...)

	return err
}
