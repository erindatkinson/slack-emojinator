IMPORTDIR ?= ./import/
EXPORTDIR ?= ./export/

## snapshot:		Builds binaries based on current code
snapshot:
	docker-compose run builder build --snapshot --clean

release:
	goreleaser release --clean

## clean:		Removes build/release/action folders
clean:
	rm -rf dist/ bin/ import/ export/

## build-docker:		Builds docker image with Pipenv requirements
build-docker: snapshot
	docker-compose build runner

## help:			Prints make target help information from comments in makefile.
help: Makefile
	@sed -n 's/^##//p' $< | sort
