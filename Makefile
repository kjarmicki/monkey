.PHONY: test repl

test:
	go test ./...

repl:
	go run main.go
