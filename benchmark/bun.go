package benchmark

import (
	"context"
	"testing"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/stdlib"
	"github.com/lauro-santana/golang-orm-benchmarks/benchmark/utils"
	"github.com/lauro-santana/golang-orm-benchmarks/model"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
)

type BunBenchmark struct {
	db  *bun.DB
	ctx context.Context
}

func NewBunBenchmark() Benchmark {
	return &BunBenchmark{ctx: context.Background()}
}

func (o *BunBenchmark) Init() error {
	config, err := pgx.ParseConfig(utils.PostgresDSN)
	if err != nil {
		panic(err)
	}
	config.DefaultQueryExecMode = pgx.QueryExecModeSimpleProtocol

	sqldb := stdlib.OpenDB(*config)
	o.db = bun.NewDB(sqldb, pgdialect.New())
	// o.db.AddQueryHook(bundebug.NewQueryHook(bundebug.WithVerbose(true)))

	return nil
}

func (o *BunBenchmark) Close() error {
	return o.db.Close()
}

func (o *BunBenchmark) Insert(b *testing.B) {
	book := model.NewBook()

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		b.StopTimer()
		book.ID = 0
		b.StartTimer()

		_, err := o.db.NewInsert().Model(book).Table("books").
			Exec(o.ctx)

		b.StopTimer()
		if err != nil {
			b.Error(err)
		}
		b.StartTimer()
	}
}

func (o *BunBenchmark) InsertBulk(b *testing.B) {
	books := model.NewBooks(utils.BulkInsertNumber)

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		b.StopTimer()
		for _, book := range books {
			book.ID = 0
		}
		b.StartTimer()

		_, err := o.db.NewInsert().Model(&books).Table("books").
			Exec(o.ctx)

		b.StopTimer()
		if err != nil {
			b.Error(err)
		}
		b.StartTimer()
	}
}

func (o *BunBenchmark) Update(b *testing.B) {
	book := model.NewBook()

	_, err := o.db.NewInsert().Model(book).Table("books").
		Exec(o.ctx)
	if err != nil {
		b.Error(err)
	}

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_, err = o.db.NewUpdate().Model(book).WherePK().Table("books").
			Exec(o.ctx)

		b.StopTimer()
		if err != nil {
			b.Error(err)
		}
		b.StartTimer()
	}
}

func (o *BunBenchmark) Delete(b *testing.B) {
	n := b.N
	books := model.NewBooks(n)

	_, err := o.db.NewInsert().Model(&books).Table("books").
		Exec(o.ctx)
	if err != nil {
		b.Error(err)
	}

	b.ReportAllocs()
	b.ResetTimer()

	var book *model.Book
	for i := 0; i < n; i++ {
		b.StopTimer()
		book = new(model.Book)
		book.ID = books[i].ID
		b.StartTimer()

		_, err = o.db.NewDelete().Model(book).WherePK().Table("books").
			Exec(o.ctx)

		b.StopTimer()
		if err != nil {
			b.Error(err)
		}
		b.StartTimer()
	}
}

func (o *BunBenchmark) FindByID(b *testing.B) {
	book := model.NewBook()
	_, err := o.db.NewInsert().Model(book).Exec(o.ctx)
	if err != nil {
		b.Error(err)
	}

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		for range utils.FindOneLoop {
			err = o.db.NewSelect().Model(book).
				Where("id = ?", book.ID).Limit(1).Scan(o.ctx)

			b.StopTimer()
			if err != nil {
				b.Error(err)
			}
			b.StartTimer()
		}
	}
}

func (o *BunBenchmark) FindPage(b *testing.B) {
	books := model.NewBooks(utils.BulkInsertPageNumber)

	_, err := o.db.NewInsert().Model(&books).Table("books").
		Exec(o.ctx)
	if err != nil {
		b.Error(err)
	}

	b.ReportAllocs()
	b.ResetTimer()

	booksPage := make([]model.Book, 0, utils.PageSize)
	for i := 0; i < b.N; i++ {
		for s := 0; s < utils.BulkInsertPageNumber; s = s + utils.PageSize {
			// ent, sqlc and goe generates the slice inside, so all makes counts
			booksPage = make([]model.Book, utils.PageSize)

			err = o.db.NewSelect().Model(&booksPage).Where("id > ?", s).Limit(utils.PageSize).Scan(o.ctx)

			b.StopTimer()
			if err != nil {
				b.Error(err)
			}
			b.StartTimer()
		}
	}
}
