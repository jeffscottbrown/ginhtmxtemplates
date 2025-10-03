.PHONY: test check-coverage coverage-report open

COVERAGE_OUT=coverage.out
COVERAGE_HTML=coverage.html

check-coverage: test
	@go tool go-test-coverage --config=./.testcoverage.yml

# Run tests with coverage, don't exit on failure
test:
	@go test -count=1 -coverpkg=./... -coverprofile=$(COVERAGE_OUT) ./... || true

# Generate HTML report from coverage profile
coverage-report: test
	@go tool cover -html=$(COVERAGE_OUT) -o $(COVERAGE_HTML)

# Open the coverage report in the default browser
open: coverage-report
	@echo "Opening $(COVERAGE_HTML)..."
	@if command -v xdg-open > /dev/null; then \
		xdg-open $(COVERAGE_HTML); \
	elif command -v open > /dev/null; then \
		open $(COVERAGE_HTML); \
	elif command -v start > /dev/null; then \
		start $(COVERAGE_HTML); \
	else \
		echo "Could not detect a command to open the browser."; \
	fi
clean:
	@rm -f $(COVERAGE_OUT) $(COVERAGE_HTML)

lint:
	@golangci-lint run ./...

check: check-coverage lint
	@echo "All checks passed!"
