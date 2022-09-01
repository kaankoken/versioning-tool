clean:
	go clean -modcache
	go mod tidy

lint:
	golint ./...

run-vet:
	go vet ./...

static-check:
	staticcheck ./...

test:
	bash ./scripts/run_tests.sh
