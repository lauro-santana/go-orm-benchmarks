package benchmark

import (
	"testing"

	"github.com/lauro-santana/golang-orm-benchmarks/benchmark/utils"
	"github.com/lauro-santana/golang-orm-benchmarks/model"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type GormBenchmark struct {
	db *gorm.DB
}

func NewGormBenchmark() Benchmark {
	return &GormBenchmark{}
}

func (o *GormBenchmark) Init() error {
	var err error
	// The config follows the performance section of the GORM documentation: https://gorm.io/docs/performance.html.
	pgConfig := postgres.New(postgres.Config{
		DSN:                  utils.PostgresDSN,
		PreferSimpleProtocol: true,
	})
	gormConfig := &gorm.Config{
		SkipDefaultTransaction: true,
		Logger:                 logger.Default.LogMode(logger.Silent),
	}
	o.db, err = gorm.Open(pgConfig, gormConfig)
	return err
}

func (o *GormBenchmark) Close() error {
	sqlDB, err := o.db.DB()
	if err != nil {
		return err
	}
	return sqlDB.Close()
}

func (o *GormBenchmark) Insert(b *testing.B) {
	book := model.NewBook()

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		b.StopTimer()
		book.ID = 0
		b.StartTimer()

		err := o.db.Create(book).Error

		b.StopTimer()
		if err != nil {
			b.Error(err)
		}
		b.StartTimer()
	}
}

func (o *GormBenchmark) InsertBulk(b *testing.B) {
	books := model.NewBooks(utils.BulkInsertNumber)

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		b.StopTimer()
		for _, book := range books {
			book.ID = 0
		}
		b.StartTimer()

		err := o.db.Create(&books).Error

		b.StopTimer()
		if err != nil {
			b.Error(err)
		}
		b.StartTimer()
	}
}

func (o *GormBenchmark) Update(b *testing.B) {
	book := model.NewBook()

	err := o.db.Create(book).Error
	if err != nil {
		b.Error(err)
	}

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		err = o.db.Save(book).Error

		b.StopTimer()
		if err != nil {
			b.Error(err)
		}
		b.StartTimer()
	}
}

func (o *GormBenchmark) Delete(b *testing.B) {
	n := b.N
	books := model.NewBooks(n)

	err := o.db.Create(books).Error
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

		err = o.db.Delete(&model.Book{}, bookID).Error

		b.StopTimer()
		if err != nil {
			b.Error(err)
		}
		b.StartTimer()
	}
}

func (o *GormBenchmark) FindByID(b *testing.B) {
	book := model.NewBook()

	err := o.db.Create(book).Error
	if err != nil {
		b.Error(err)
	}

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		for range utils.FindOneLoop {
			err = o.db.First(book, book.ID).Error

			b.StopTimer()
			if err != nil {
				b.Error(err)
			}
			b.StartTimer()
		}
	}
}

func (o *GormBenchmark) FindPage(b *testing.B) {
	books := model.NewBooks(utils.BulkInsertPageNumber)

	// Persist it in batches > https://github.com/bitmagnet-io/bitmagnet/issues/126.
	batches := model.Chunk(books, utils.BatchSize)
	for _, chunk := range batches {
		err := o.db.Create(chunk).Error
		if err != nil {
			b.Error(err)
		}

	}

	b.ReportAllocs()
	b.ResetTimer()

	booksPage := make([]model.Book, 0, utils.PageSize)
	for i := 0; i < b.N; i++ {
		for s := 0; s < utils.BulkInsertPageNumber; s = s + utils.PageSize {
			// ent, sqlc and goe generates the slice inside, so all makes counts
			booksPage = make([]model.Book, utils.PageSize)

			err := o.db.Limit(utils.PageSize).Where("id > ?", s).Find(&booksPage).Error

			b.StopTimer()
			if err != nil {
				b.Error(err)
			}
			b.StartTimer()
		}
	}
}
