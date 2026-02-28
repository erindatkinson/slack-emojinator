IMPORTDIR ?= ./import/
EXPORTDIR ?= ./export/

## snapshot:		Builds binaries based on current code
snapshot: bindata
	goreleaser build --snapshot --clean

release: bindata
	goreleaser release --clean

bindata:
	go-bindata -modtime 1771972135 -pkg templates -o ./internal/templates/bindata.go templates/*

## clean:		Removes build/release/action folders
clean:
	rm -rf dist/ bin/ import/ export/

test:
	go test ./...

coverage:
	go test --coverprofile .coverage
	go tool cover -html=.coverage

## help:			Prints make target help information from comments in makefile.
help: Makefile
	@sed -n 's/^##//p' $< | sort
