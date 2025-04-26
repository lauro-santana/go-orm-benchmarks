package benchmark

import (
	"context"
	"testing"

	"github.com/lauro-santana/golang-orm-benchmarks/benchmark/utils"
	"github.com/lauro-santana/golang-orm-benchmarks/model"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

var columns = []string{"isbn", "title", "author", "genre", "quantity", "publicized_at"}

type PgxBenchmark struct {
	db  *pgxpool.Pool
	ctx context.Context
}

func NewPgxBenchmark() Benchmark {
	return &PgxBenchmark{
		ctx: context.Background(),
	}
}

func (p *PgxBenchmark) Init() error {
	var err error
	p.db, err = pgxpool.New(p.ctx, utils.PostgresDSN)
	return err
}

func (p *PgxBenchmark) Close() error {
	p.db.Close()
	return nil
}

func (p *PgxBenchmark) Insert(b *testing.B) {
	book := model.NewBook()

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_, err := p.db.Exec(p.ctx, utils.InsertQuery,
			book.ISBN, book.Title, book.Author, book.Genre, book.Quantity, book.PublicizedAt)

		b.StopTimer()
		if err != nil {
			b.Error(err)
		}
		b.StartTimer()
	}
}

func (p *PgxBenchmark) InsertBulk(b *testing.B) {
	var rows = make([][]interface{}, 0)
	for _, book := range model.NewBooks(utils.BulkInsertNumber) {
		rows = append(rows, []interface{}{book.ISBN, book.Title, book.Author, book.Genre, book.Quantity, book.PublicizedAt})
	}

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_, err := p.db.CopyFrom(p.ctx, pgx.Identifier{"books"}, columns, pgx.CopyFromRows(rows))

		if err != nil {
			b.Error(err)
		}
	}
}

func (p *PgxBenchmark) Update(b *testing.B) {
	book := model.NewBook()
	var id int64
	err := p.db.QueryRow(p.ctx, utils.InsertReturningIDQuery,
		book.ISBN, book.Title, book.Author, book.Genre, book.Quantity, book.PublicizedAt).Scan(&id)
	if err != nil {
		b.Error(err)
	}

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_, err = p.db.Exec(p.ctx, utils.UpdateQuery,
			book.ISBN, book.Title, book.Author, book.Genre, book.Quantity, book.PublicizedAt, id)

		if err != nil {
			b.Error(err)
		}
	}
}

func (p *PgxBenchmark) Delete(b *testing.B) {
	book := model.NewBook()
	savedIDs := make([]int64, b.N)
	for i := 0; i < b.N; i++ {
		var id int64
		err := p.db.QueryRow(p.ctx, utils.InsertReturningIDQuery,
			book.ISBN, book.Title, book.Author, book.Genre, book.Quantity, book.PublicizedAt).Scan(&id)
		if err != nil {
			b.Error(err)
		}
		savedIDs[i] = id
	}

	b.ReportAllocs()
	b.ResetTimer()

	var bookID int64
	for i := 0; i < b.N; i++ {
		b.StopTimer()
		bookID = savedIDs[i]
		b.StartTimer()

		_, err := p.db.Exec(p.ctx, utils.DeleteQuery, bookID)

		if err != nil {
			b.Error(err)
		}
	}
}

func (p *PgxBenchmark) FindByID(b *testing.B) {
	book := model.NewBook()
	var id int64
	err := p.db.QueryRow(p.ctx, utils.InsertReturningIDQuery,
		book.ISBN, book.Title, book.Author, book.Genre, book.Quantity, book.PublicizedAt).Scan(&id)
	if err != nil {
		b.Error(err)
	}

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		for range utils.FindOneLoop {
			var foundBook model.Book
			err := p.db.QueryRow(p.ctx, utils.SelectByIDQuery, id).Scan(
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

func (p *PgxBenchmark) FindPage(b *testing.B) {
	var rows = make([][]interface{}, 0)
	for _, book := range model.NewBooks(utils.BulkInsertPageNumber) {
		rows = append(rows, []interface{}{book.ISBN, book.Title, book.Author, book.Genre, book.Quantity, book.PublicizedAt})
	}

	_, err := p.db.CopyFrom(p.ctx, pgx.Identifier{"books"}, columns, pgx.CopyFromRows(rows))
	if err != nil {
		b.Error(err)
	}

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		for s := 0; s < utils.BulkInsertPageNumber; s = s + utils.PageSize {
			// making slices will count on raw benchmarks
			booksPage = make([]model.Book, 0, utils.PageSize)

			result, err := p.db.Query(p.ctx, utils.SelectPaginatingQuery, s, utils.PageSize)

			// checking the error will count on raw benchmarks
			if err != nil {
				b.Error(err)
			}

			for result.Next() {
				var book model.Book
				if err = result.Scan(
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
