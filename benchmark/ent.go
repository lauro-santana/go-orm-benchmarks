package benchmark

import (
	"context"
	"database/sql"
	"testing"

	"github.com/lauro-santana/golang-orm-benchmarks/benchmark/ent"
	"github.com/lauro-santana/golang-orm-benchmarks/benchmark/ent/book"
	"github.com/lauro-santana/golang-orm-benchmarks/benchmark/utils"
	"github.com/lauro-santana/golang-orm-benchmarks/model"

	"entgo.io/ent/dialect"
	entsql "entgo.io/ent/dialect/sql"

	// Postgres driver.
	_ "github.com/jackc/pgx/v5/stdlib"
)

type EntBenchmark struct {
	db  *ent.Client
	ctx context.Context
}

func NewEntBenchmark() Benchmark {
	return &EntBenchmark{ctx: context.Background()}
}

func (o *EntBenchmark) Init() error {
	db, err := sql.Open("pgx", utils.PostgresDSN)
	if err != nil {
		return err
	}
	drv := entsql.OpenDB(dialect.Postgres, db)
	o.db = ent.NewClient(ent.Driver(drv))
	return nil
}

func (o *EntBenchmark) Close() error {
	return o.db.Close()
}

func (o *EntBenchmark) Insert(b *testing.B) {
	newBook := model.NewBook()

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		b.StopTimer()
		newBook.ID = 0
		b.StartTimer()

		_, err := o.db.Book.
			Create().
			SetIsbn(newBook.ISBN).
			SetTitle(newBook.Title).
			SetAuthor(newBook.Author).
			SetGenre(newBook.Genre).
			SetQuantity(newBook.Quantity).
			SetPublicizedAt(newBook.PublicizedAt).
			Save(o.ctx)

		b.StopTimer()
		if err != nil {
			b.Error(err)
		}
		b.StartTimer()
	}
}

func (o *EntBenchmark) InsertBulk(b *testing.B) {
	books := model.NewBooks(utils.BulkInsertNumber)

	b.ReportAllocs()
	b.ResetTimer()

	batch := make([]*ent.BookCreate, len(books))
	for i, newBook := range books {
		batch[i] = o.db.Book.Create().
			SetIsbn(newBook.ISBN).
			SetTitle(newBook.Title).
			SetAuthor(newBook.Author).
			SetGenre(newBook.Genre).
			SetQuantity(newBook.Quantity).
			SetPublicizedAt(newBook.PublicizedAt)
	}

	for i := 0; i < b.N; i++ {
		_, err := o.db.Book.CreateBulk(batch...).Save(o.ctx)

		b.StopTimer()
		if err != nil {
			b.Error(err)
		}
		b.StartTimer()
	}

}

func (o *EntBenchmark) Update(b *testing.B) {
	newBook := model.NewBook()

	saved, err := o.db.Book.
		Create().
		SetIsbn(newBook.ISBN).
		SetTitle(newBook.Title).
		SetAuthor(newBook.Author).
		SetGenre(newBook.Genre).
		SetQuantity(newBook.Quantity).
		SetPublicizedAt(newBook.PublicizedAt).
		Save(o.ctx)
	if err != nil {
		b.Error(err)
	}

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_, err = o.db.Book.
			UpdateOneID(saved.ID).
			SetIsbn(newBook.ISBN).
			SetTitle(newBook.Title).
			SetAuthor(newBook.Author).
			SetGenre(newBook.Genre).
			SetQuantity(newBook.Quantity).
			SetPublicizedAt(newBook.PublicizedAt).
			Save(o.ctx)

		b.StopTimer()
		if err != nil {
			b.Error(err)
		}
		b.StartTimer()
	}
}

func (o *EntBenchmark) Delete(b *testing.B) {
	n := b.N
	books := model.NewBooks(n)
	batch := make([]*ent.BookCreate, len(books))
	for i, newBook := range books {
		batch[i] = o.db.Book.Create().
			SetIsbn(newBook.ISBN).
			SetTitle(newBook.Title).
			SetAuthor(newBook.Author).
			SetGenre(newBook.Genre).
			SetQuantity(newBook.Quantity).
			SetPublicizedAt(newBook.PublicizedAt)
	}

	saved, err := o.db.Book.CreateBulk(batch...).Save(o.ctx)
	if err != nil {
		b.Error(err)
	}

	b.ReportAllocs()
	b.ResetTimer()

	var bookID int
	for i := 0; i < n; i++ {
		b.StopTimer()
		bookID = saved[i].ID
		b.StartTimer()

		err = o.db.Book.
			DeleteOneID(bookID).
			Exec(o.ctx)

		b.StopTimer()
		if err != nil {
			b.Error(err)
		}
		b.StartTimer()
	}
}

func (o *EntBenchmark) FindByID(b *testing.B) {
	newBook := model.NewBook()
	book, err := o.db.Book.
		Create().
		SetIsbn(newBook.ISBN).
		SetTitle(newBook.Title).
		SetAuthor(newBook.Author).
		SetGenre(newBook.Genre).
		SetQuantity(newBook.Quantity).
		SetPublicizedAt(newBook.PublicizedAt).
		Save(o.ctx)
	if err != nil {
		b.Error(err)
	}

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		for range utils.FindOneLoop {
			_, err = o.db.Book.Get(o.ctx, book.ID)

			b.StopTimer()
			if err != nil {
				b.Error(err)
			}
			b.StartTimer()
		}
	}
}

func (o *EntBenchmark) FindPage(b *testing.B) {
	books := model.NewBooks(utils.BulkInsertPageNumber)
	batch := make([]*ent.BookCreate, len(books))
	for i, newBook := range books {
		batch[i] = o.db.Book.Create().
			SetIsbn(newBook.ISBN).
			SetTitle(newBook.Title).
			SetAuthor(newBook.Author).
			SetGenre(newBook.Genre).
			SetQuantity(newBook.Quantity).
			SetPublicizedAt(newBook.PublicizedAt)
	}
	_, err := o.db.Book.CreateBulk(batch...).Save(o.ctx)
	if err != nil {
		b.Error(err)
	}

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		for s := 0; s < utils.BulkInsertPageNumber; s = s + utils.PageSize {
			_, err = o.db.Book.
				Query().
				Where(book.IDGT(s)).
				Limit(utils.PageSize).
				All(o.ctx)

			b.StopTimer()
			if err != nil {
				b.Error(err)
			}
			b.StartTimer()
		}
	}
}
