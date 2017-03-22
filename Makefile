default: build

build:
	go build ./cmd/exchangerates
	go build ./cmd/exchangeratesd

test:
	go test $$(go list ./... | grep -v vendor)

.PHONY: build test
