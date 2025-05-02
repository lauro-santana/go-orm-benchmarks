## Golang ORM Benchmarks

Must have Go 1.23+, Make and Docker installed.

#### This repository contains benchmarks for the following projects

- [GORM](https://gorm.io/)
- [Ent](https://entgo.io/)
- [Bun](https://bun.uptrace.dev/)
- [Sqlc](https://sqlc.dev/)
- [GOE](https://github.com/go-goe/goe)

#### And also, pure SQL benchmarks using

- [pgx](https://github.com/jackc/pgx)
- [database/sql](https://pkg.go.dev/database/sql)

Run all benchmarks using the following command:

```bash
make benchmark-all
```

<p>If you want to run a specific benchmark, you can use the following commands:

```bash
make benchmark-insert
make benchmark-insert-bulk
make benchmark-update
make benchmark-delete
make benchmark-select-one
make benchmark-select-page
```

Modeling credits: [efectn/go-orm-benchmarks](https://github.com/efectn/go-orm-benchmarks) and [andreiac-silva/golang-orm-benchmarks](https://github.com/andreiac-silva/golang-orm-benchmarks).
