$PHONY: clean test build

clean: 
		rm cholibot || true

test:
	go test ./...
		
build:
		CGO_ENABLED=0 GOOS=linux go build -o cholibot ./cmd/monolith/main.go
		sha256sum cholibot > cholibot.sha256

