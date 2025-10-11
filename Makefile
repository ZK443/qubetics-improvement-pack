.PHONY: build vet test cover

build:
	@echo "==> go build ./..."
	go build ./... 2>&1 | tee build.log

vet:
	@echo "==> go vet ./..."
	go vet ./...

test:
	@echo "==> go test -v ./..."
	go test -v ./...

cover:
	@echo "==> go test -coverprofile=coverage.out ./..."
	go test -coverprofile=coverage.out ./...
	@echo "==> convert to HTML"
	go tool cover -html=coverage.out -o coverage.html || true
