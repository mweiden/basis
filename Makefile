test:
	go test -v ./... && python -m unittest discover -s ./ -p '*_test.py'
	@echo "Test result $0"

test-only:
	find . -name '*.go' | entr go test -v $(SUBDIR)*.go

format:
	go fmt ./...
