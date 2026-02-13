IMPORTDIR ?= ./import/
EXPORTDIR ?= ./export/

## snapshot:		Builds binaries based on current code
snapshot: bindata
	docker-compose run builder build --snapshot --clean

release: bindata
	goreleaser release --clean

bindata:
	go-bindata -pkg templates -o ./internal/templates/bindata.go templates/*

## clean:		Removes build/release/action folders
clean:
	rm -rf dist/ bin/ import/ export/

## build-docker:		Builds docker image with Pipenv requirements
build-docker: snapshot
	docker-compose build runner

test:
	go test ./...

coverage:
	go test --coverprofile .coverage
	go tool cover -html=.coverage

## help:			Prints make target help information from comments in makefile.
help: Makefile
	@sed -n 's/^##//p' $< | sort
