package main

import (
	"flag"
	"fmt"
	"log"
	"math/rand"
	"os"
	"slices"
	"testing"
	"text/tabwriter"
	"time"

	"github.com/lauro-santana/golang-orm-benchmarks/benchmark"

	// Auto load .env file.
	_ "github.com/joho/godotenv/autoload"
)

const (
	all = "all"

	insertOp     = "insert"
	insertBulkOp = "insert-bulk"
	updateOp     = "update"
	deleteOp     = "delete"
	selectOne    = "select-one"
	selectPage   = "select-page"

	raw  = "database/sql"
	pgx  = "pgx"
	bun  = "bun"
	gorm = "gorm"
	ent  = "ent"
	sqlc = "sqlc"
	goe  = "goe"
)

var (
	benchmarksMap   = map[string]benchmark.Benchmark{}
	validOperations = []string{insertOp, insertBulkOp, updateOp, deleteOp, selectOne, selectPage}
)

func main() {
	operation := flag.String("operation", selectOne, "Specify the operation to run")
	flag.Parse()

	if operation == nil && *operation != all && slices.Contains(validOperations, *operation) {
		log.Fatal("define a valid orm or operation")
	}

	loadBenchmarks()
	shuffleBenchmarksMap()
	results := executeBenchmarks(*operation)
	printBenchmark(results, *operation)
}

func loadBenchmarks() {
	benchmarksMap[raw] = benchmark.NewRawBenchmark()
	benchmarksMap[pgx] = benchmark.NewPgxBenchmark()
	benchmarksMap[bun] = benchmark.NewBunBenchmark()
	benchmarksMap[gorm] = benchmark.NewGormBenchmark()
	benchmarksMap[ent] = benchmark.NewEntBenchmark()
	benchmarksMap[sqlc] = benchmark.NewSqlcBenchmark()
	benchmarksMap[goe] = benchmark.NewGoeBenchmark()
}

func shuffleBenchmarksMap() {
	source := rand.NewSource(time.Now().UnixNano())
	rng := rand.New(source)
	keys := make([]string, 0, len(benchmarksMap))
	for key := range benchmarksMap {
		keys = append(keys, key)
	}
	rng.Shuffle(len(keys), func(i, j int) {
		keys[i], keys[j] = keys[j], keys[i]
	})
	shuffledMap := make(map[string]benchmark.Benchmark)
	for _, key := range keys {
		shuffledMap[key] = benchmarksMap[key]
	}
	benchmarksMap = shuffledMap
}

func executeBenchmarks(operation string) []benchmark.ResultWrapper {
	var results []benchmark.ResultWrapper
	for ormName, b := range benchmarksMap {
		results = append(results, doExecuteBenchmarks(b, ormName, operation))
	}
	return results
}

func doExecuteBenchmarks(b benchmark.Benchmark, orm, operation string) benchmark.ResultWrapper {
	benchmark.BeforeBenchmark()
	wrapper := benchmark.ResultWrapper{}
	wrapper.Orm = orm
	err := b.Init()
	if err != nil {
		wrapper.Err = err
	}
	resultMap := make(map[string]testing.BenchmarkResult)
	operations := map[string]func(*testing.B){
		insertOp:     b.Insert,
		insertBulkOp: b.InsertBulk,
		updateOp:     b.Update,
		deleteOp:     b.Delete,
		selectOne:    b.FindByID,
		selectPage:   b.FindPage,
	}
	if operation == all {
		for op, f := range operations {
			resultMap[op] = testing.Benchmark(f)
		}
		wrapper.Benchmarks = resultMap
		return wrapper
	}
	wrapper.Benchmarks = map[string]testing.BenchmarkResult{
		operation: testing.Benchmark(operations[operation]),
	}
	return wrapper
}

func printBenchmark(results []benchmark.ResultWrapper, operation string) {
	table := new(tabwriter.Writer)
	table.Init(os.Stdout, 0, 8, 2, '\t', tabwriter.AlignRight)
	if operation == all {
		doPrintBenchmark(table, results, validOperations...)
	} else {
		doPrintBenchmark(table, results, operation)
	}
}

func doPrintBenchmark(table *tabwriter.Writer, results []benchmark.ResultWrapper, operations ...string) {
	for _, op := range operations {
		_, _ = fmt.Fprint(table, "\n")
		_, _ = fmt.Fprintf(table, "Operation: %s\n", op)

		for _, r := range results {
			result, ok := r.Benchmarks[op]
			if !ok {
				continue
			}
			_, _ = fmt.Fprintf(table, "%s:\t%d\t%d ns/op\t%d B/op\t%d allocs/op\n",
				r.Orm,
				result.N,
				result.NsPerOp(),
				result.AllocedBytesPerOp(),
				result.AllocsPerOp(),
			)
		}

		_ = table.Flush()
	}
}
