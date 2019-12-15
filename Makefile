build:
	go generate
	go build -mod=vendor

.PHONY: build
