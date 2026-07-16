.PHONY: run-tests

run-tests:
	go test -failfast -cover ./...
