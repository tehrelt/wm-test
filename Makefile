APP_NAME=wm-test
APP_BIN=bin/$(APP_NAME)

.PHONY: build run
build:
	go build -o $(APP_BIN) ./cmd/app/main.go

run: build
	./$(APP_BIN)
