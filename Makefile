run: 
	go run main.go
help:
	go run main.go --help
format:
	gofmt -s -w .
.PHONY: run format help