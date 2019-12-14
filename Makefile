build: bindata.go
	go build

bindata.go:
	go run ./tools/go-bindata/ -prefix data data/...

.PHONY: build