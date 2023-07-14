.PHONY: test
test: ## runs unit tests with coverage
	@go test -v -failfast `go list ./... | grep -v "mocks"` -failfast --coverprofile="coverage.tmp.out" -covermode=count;
	@cat coverage.tmp.out | grep -v "_mock.go" > coverage.out;
	@go tool cover -func coverage.out | grep total | awk '{print $3}';
	@go tool cover -html=coverage.out -o coverage.html;