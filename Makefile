# Run all benchmarks sequentially and save to results/ with timestamps
benchmark-all: benchmark-insert benchmark-insert-bulk benchmark-update \
              benchmark-delete benchmark-select-one benchmark-select-page
	@echo "All benchmarks completed. Results saved to results/"

# Create results directory if not exists
results:
	mkdir -p results

# Generic benchmark runner with output redirection
benchmark-insert: results
	docker compose up -d --no-recreate
	go run main.go -operation insert | tee results/insert-$(shell date +%Y%m%d-%H%M%S).log

benchmark-insert-bulk: results
	docker compose up -d --no-recreate
	go run main.go -operation insert-bulk | tee results/insert-bulk-$(shell date +%Y%m%d-%H%M%S).log

benchmark-update: results
	docker compose up -d --no-recreate
	go run main.go -operation update | tee results/update-$(shell date +%Y%m%d-%H%M%S).log

benchmark-delete: results
	docker compose up -d --no-recreate
	go run main.go -operation delete | tee results/delete-$(shell date +%Y%m%d-%H%M%S).log

benchmark-select-one: results
	docker compose up -d --no-recreate
	go run main.go -operation select-one | tee results/select-one-$(shell date +%Y%m%d-%H%M%S).log

benchmark-select-page: results
	docker compose up -d --no-recreate
	go run main.go -operation select-page | tee results/select-page-$(shell date +%Y%m%d-%H%M%S).log
