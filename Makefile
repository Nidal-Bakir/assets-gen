.DEFAULT_GOAL := build

build:
	@echo "Building..."
	@go build -o main cmd/cli/main.go


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
	@rm -f main
	@rm -rf assets_gen_out/
	@echo "Done"
