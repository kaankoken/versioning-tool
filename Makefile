build:
	go build main.go

clean:
	go clean -modcache
	go mod tidy

lint:
	bash ./scripts/run_command.sh "golint ./..."	

run-vet:
	bash ./scripts/run_command.sh "go vet ./..."

static-check:
	bash ./scripts/run_command.sh "staticcheck ./..."

test:
	bash ./scripts/run_tests.sh
