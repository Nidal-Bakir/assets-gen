.DEFAULT_GOAL := build

.PHONY: build test clean

build:
	@echo "Building..."
	@go build -o build/assets-gen cmd/assetsgen/main.go

run:
	@go run cmd/cli/main.go

test:
	@echo "Testing..."
	@go test ./... -v

chk:
	@echo "Runing go vet & staticcheck..."
	@go vet ./...
	@staticcheck ./...

fmt:
	@echo "Formating..."
	@go fmt ./...

clean:
	@echo "Cleaning..."
	@rm -rf build/
	@rm -rf assets_gen_out/
	@echo "Done"
