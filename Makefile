.PHONY: start_server

.DEFAULT_GOAL := start_server

start_server:
	go run ./cmd/server
