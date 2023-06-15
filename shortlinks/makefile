run:
	@echo "Running..."
	@go run main.go
	@echo "Done."

build:
	@echo "Building..."
	@go build -o bin/main main.go
	@echo "Done."

clean:
	@echo "Cleaning..."
	@rm -rf bin
	@echo "Done."

test:
	@echo "Testing..."
	@go test ./...
	@echo "Done."

fmt:
	@echo "Formatting..."
	@trunk fmt
	@echo "Done."

check:
	@echo "Checking..."
	@trunk check
	@echo "Done."

fmt_all:
	@echo "Formatting all..."
	@trunk fmt --all
	@echo "Done."

check_all:
	@echo "Checking all..."
	@trunk check --all
	@echo "Done."
