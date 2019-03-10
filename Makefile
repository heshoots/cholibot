$PHONY: clean test build

clean: 
		rm cholibot || true

test:
	go test ./...
		
build:
		CGO_ENABLED=0 GOOS=linux go build -o cholibot ./cmd/cholibot/main.go

