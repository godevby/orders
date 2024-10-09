SHELL = /bin/bash

run:
	go run cmd/orders/main.go

# ----------------------------------------------------------------------
# Modules support

tidy:
	go mod tidy
	go mod vendor