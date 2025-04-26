package benchmark

import (
	"context"
	"testing"

	"github.com/lauro-santana/golang-orm-benchmarks/benchmark/sqlc/repository"
	"github.com/lauro-santana/golang-orm-benchmarks/benchmark/utils"
	"github.com/lauro-santana/golang-orm-benchmarks/model"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
)

type SqlcBenchmark struct {
	repository *repository.Queries
	db         *pgxpool.Pool
	ctx        context.Context
}

func NewSqlcBenchmark() Benchmark {
	return &SqlcBenchmark{ctx: context.Background()}
}

func (s *SqlcBenchmark) Init() error {
	conn, err := pgxpool.New(context.Background(), utils.PostgresDSN)
	if err != nil {
		return err
	}
	s.db = conn
	s.repository = repository.New(conn)
	return nil
}

func (s *SqlcBenchmark) Close() error {
	s.db.Close()
	return nil
}

func (s *SqlcBenchmark) Insert(b *testing.B) {
	book := model.NewBook()

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		b.StopTimer()
		book.ID = 0
		b.StartTimer()

		err := s.repository.Create(s.ctx, repository.CreateParams{
			Isbn:         book.ISBN,
			Title:        book.Title,
			Author:       book.Author,
			Genre:        book.Genre,
			Quantity:     int32(book.Quantity),
			PublicizedAt: pgtype.Timestamp{Time: book.PublicizedAt, Valid: true},
		})

		b.StopTimer()
		if err != nil {
			b.Error(err)
		}
		b.StartTimer()
	}
}

func (s *SqlcBenchmark) InsertBulk(b *testing.B) {
	books := model.NewBooks(utils.BulkInsertNumber)

	b.ReportAllocs()
	b.ResetTimer()

	batch := make([]repository.CreateManyParams, len(books))
	for i, newBook := range books {
		batch[i] = repository.CreateManyParams{
			Isbn:         newBook.ISBN,
			Title:        newBook.Title,
			Author:       newBook.Title,
			Genre:        newBook.Genre,
			Quantity:     int32(newBook.Quantity),
			PublicizedAt: pgtype.Timestamp{Time: newBook.PublicizedAt, Valid: true},
		}
	}

	for i := 0; i < b.N; i++ {
		_, err := s.repository.CreateMany(s.ctx, batch)

		if err != nil {
			b.Error(err)
		}
	}
}

func (s *SqlcBenchmark) Update(b *testing.B) {
	book := model.NewBook()

	id, err := s.repository.CreateReturningID(s.ctx, repository.CreateReturningIDParams{
		Isbn:         book.ISBN,
		Title:        book.Title,
		Author:       book.Author,
		Genre:        book.Genre,
		Quantity:     int32(book.Quantity),
		PublicizedAt: pgtype.Timestamp{Time: book.PublicizedAt, Valid: true},
	})
	if err != nil {
		b.Error(err)
	}

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		err = s.repository.Update(s.ctx, repository.UpdateParams{
			ID:           id,
			Isbn:         book.ISBN,
			Title:        book.Title,
			Author:       book.Author,
			Genre:        book.Genre,
			Quantity:     int32(book.Quantity),
			PublicizedAt: pgtype.Timestamp{Time: book.PublicizedAt, Valid: true},
		})

		if err != nil {
			b.Error(err)
		}
	}
}

func (s *SqlcBenchmark) Delete(b *testing.B) {
	n := b.N
	book := model.NewBook()
	bookIDs := make([]int32, n)
	for i := 0; i < n; i++ {
		id, err := s.repository.CreateReturningID(s.ctx, repository.CreateReturningIDParams{
			Isbn:         book.ISBN,
			Title:        book.Title,
			Author:       book.Author,
			Genre:        book.Genre,
			Quantity:     int32(book.Quantity),
			PublicizedAt: pgtype.Timestamp{Time: book.PublicizedAt, Valid: true},
		})
		if err != nil {
			b.Error(err)
		}
		bookIDs[i] = id
	}

	b.ReportAllocs()
	b.ResetTimer()

	var bookID int32
	for i := 0; i < b.N; i++ {
		b.StopTimer()
		bookID = bookIDs[i]
		b.StartTimer()

		err := s.repository.Delete(s.ctx, bookID)

		if err != nil {
			b.Error(err)
		}
	}
}

func (s *SqlcBenchmark) FindByID(b *testing.B) {
	book := model.NewBook()
	id, err := s.repository.CreateReturningID(s.ctx, repository.CreateReturningIDParams{
		Isbn:         book.ISBN,
		Title:        book.Title,
		Author:       book.Author,
		Genre:        book.Genre,
		Quantity:     int32(book.Quantity),
		PublicizedAt: pgtype.Timestamp{Time: book.PublicizedAt, Valid: true},
	})
	if err != nil {
		b.Error(err)
	}

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		for range utils.FindOneLoop {
			_, err := s.repository.Get(s.ctx, id)

			// Get not checking the error, so will count on benchmark
			if err != nil {
				b.Error(err)
			}
		}
	}
}

func (s *SqlcBenchmark) FindPage(b *testing.B) {
	book := model.NewBook()
	bookIDs := make([]int32, utils.BulkInsertPageNumber)
	for i := 0; i < utils.BulkInsertPageNumber; i++ {
		id, err := s.repository.CreateReturningID(s.ctx, repository.CreateReturningIDParams{
			Isbn:         book.ISBN,
			Title:        book.Title,
			Author:       book.Author,
			Genre:        book.Genre,
			Quantity:     int32(book.Quantity),
			PublicizedAt: pgtype.Timestamp{Time: book.PublicizedAt, Valid: true},
		})
		if err != nil {
			b.Error(err)
		}
		bookIDs[i] = id
	}

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		for size := 0; size < utils.BulkInsertPageNumber; size = size + utils.PageSize {
			_, err := s.repository.ListPaginating(s.ctx, repository.ListPaginatingParams{
				ID:    int32(size),
				Limit: utils.PageSize,
			})

			b.StopTimer()
			if err != nil {
				b.Error(err)
			}
			b.StartTimer()
		}
	}
}
