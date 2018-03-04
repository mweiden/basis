build:
	go build -v mattcoin.go

test:
	go test -v ./...
	@echo "Test result $0"

test-only:
	find . -name '*.go' | entr go test -v $(SUBDIR)*.go

format:
	go fmt ./...
