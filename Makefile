test:
	go test -v ./... && python -m unittest discover -s ./ -p '*_test.py'
	@echo "Test result $0"
