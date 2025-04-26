package benchmark

import (
	"testing"

	"github.com/go-goe/goe"
	"github.com/go-goe/goe/query/where"
	"github.com/go-goe/postgres"
	"github.com/lauro-santana/golang-orm-benchmarks/benchmark/utils"
	"github.com/lauro-santana/golang-orm-benchmarks/model"
)

type Database struct {
	Book *model.Book
	*goe.DB
}

type GoeBenchmark struct {
	db *Database
}

func NewGoeBenchmark() Benchmark {
	return &GoeBenchmark{}
}

func (o *GoeBenchmark) Init() (err error) {
	o.db, err = goe.Open[Database](postgres.Open(utils.PostgresDSN, postgres.Config{}))
	return err
}

func (o *GoeBenchmark) Close() error {
	return goe.Close(o.db)
}

func (o *GoeBenchmark) Insert(b *testing.B) {
	book := model.NewBook()

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		b.StopTimer()
		book.ID = 0
		b.StartTimer()

		err := goe.Insert(o.db.Book).One(book)

		b.StopTimer()
		if err != nil {
			b.Error(err)
		}
		b.StartTimer()
	}
}

func (o *GoeBenchmark) InsertBulk(b *testing.B) {
	books := model.NewBooksNoPtr(utils.BulkInsertNumber)

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		b.StopTimer()
		for i := range books {
			books[i].ID = 0
		}
		b.StartTimer()

		err := goe.Insert(o.db.Book).All(books)

		b.StopTimer()
		if err != nil {
			b.Error(err)
		}
		b.StartTimer()
	}
}

func (o *GoeBenchmark) Update(b *testing.B) {
	book := model.NewBook()

	err := goe.Insert(o.db.Book).One(book)
	if err != nil {
		b.Error(err)
	}

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		err = goe.Save(o.db.Book).ByValue(*book)

		b.StopTimer()
		if err != nil {
			b.Error(err)
		}
		b.StartTimer()
	}
}

func (o *GoeBenchmark) Delete(b *testing.B) {
	n := b.N
	books := model.NewBooksNoPtr(n)

	err := goe.Insert(o.db.Book).All(books)
	if err != nil {
		b.Error(err)
	}

	b.ReportAllocs()
	b.ResetTimer()

	var bookID int64
	for i := 0; i < b.N; i++ {
		b.StopTimer()
		bookID = books[i].ID
		b.StartTimer()

		err = goe.Delete(o.db.Book).Where(where.Equals(&o.db.Book.ID, bookID))

		b.StopTimer()
		if err != nil {
			b.Error(err)
		}
		b.StartTimer()
	}
}

func (o *GoeBenchmark) FindByID(b *testing.B) {
	book := model.NewBook()

	err := goe.Insert(o.db.Book).One(book)
	if err != nil {
		b.Error(err)
	}

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		for range utils.FindOneLoop {
			_, err = goe.Find(o.db.Book).ById(model.Book{ID: book.ID})

			b.StopTimer()
			if err != nil {
				b.Error(err)
			}
			b.StartTimer()
		}
	}
}

func (o *GoeBenchmark) FindPage(b *testing.B) {
	books := model.NewBooksNoPtr(utils.BulkInsertPageNumber)

	err := goe.Insert(o.db.Book).All(books)
	if err != nil {
		b.Error(err)
	}

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		for s := int64(0); s < utils.BulkInsertPageNumber; s = s + utils.PageSize {
			_, err = goe.Select(o.db.Book).From(o.db.Book).Take(utils.PageSize).Where(where.Greater(&o.db.Book.ID, s)).AsSlice()

			b.StopTimer()
			if err != nil {
				b.Error(err)
			}
			b.StartTimer()
		}
	}
}
