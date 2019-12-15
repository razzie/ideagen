build: bindata.go
	go generate
	go build

.PHONY: build