.PHONY: test testacc

test: deps
	go test -v ./...

testacc: deps
	TESTACC=1 go test -p 1 -v ./... -run="TestAcc"

deps:
	go mod download
	go mod tidy
