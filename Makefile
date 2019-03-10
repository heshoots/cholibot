$PHONY: clean

clean: 
		rm cholibot

cholibot:
		CGO_ENABLED=0 GOOS=linux go build -o cholibot ./cmd/cholibot/main.go

build: cholibot

