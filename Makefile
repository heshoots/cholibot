$PHONY: clean test build

clean: 
		rm cholibot || true

test:
	go test ./...
		
build:
		CGO_ENABLED=0 GOOS=linux go build -ldflags "-X main.version=$(TRAVIS_TAG)" -o cholibot ./cmd/monolith/main.go
		sha256sum cholibot > cholibot.sha256

