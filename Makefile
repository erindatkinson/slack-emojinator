IMPORTDIR ?= ./import/
EXPORTDIR ?= ./export/

## all:			Sets up pipenv requirements both locally and for docker
all: setup build-docker

## lint:			Runs pylint on the directory
lint:
	pipenv run pylint ./

## build-docker:		Builds docker image with Pipenv requirements
build-docker:
	docker-compose build

## setup:			Installs Pipenv requirements locally
setup:
	pipenv install

## docker-export:		Exports emoji from your slack team to the ./export/ directory using docker
docker-export: 
	docker-compose run emoji-export

## docker-import:		Imports emoji from the ./import/ directory to your slack team using docker
docker-import: 
	docker-compose run emoji-import

## download:		Exports emoji from your slack team to the ./export/ directory
download: 
	pipenv run python main.py export $(EXPORTDIR)

## upload:		Imports emoji from the ./import/ directory to your slack team
upload: 
	pipenv run python main.py import $(IMPORTDIR)

stats: 
	pipenv run python main.py stats

## help:			Prints make target help information from comments in makefile.
help: Makefile
	@sed -n 's/^##//p' $< | sort
