benchmark-insert: # Run insert benchmarks
	docker compose up -d --no-recreate
	go run main.go -operation insert

benchmark-insert-bulk: # Run insert bulk benchmarks
	docker compose up -d --no-recreate
	go run main.go -operation insert-bulk

benchmark-update: # Run update benchmarks
	docker compose up -d --no-recreate
	go run main.go -operation update

benchmark-delete: # Run delete benchmarks
	docker compose up -d --no-recreate
	go run main.go -operation delete

benchmark-select-one: # Run select one benchmarks
	docker compose up -d --no-recreate
	go run main.go -operation select-one

benchmark-select-page: # Run select page benchmarks
	docker compose up -d --no-recreate
	go run main.go -operation select-page
