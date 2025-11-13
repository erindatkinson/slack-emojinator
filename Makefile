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
docker-export: needs-envs
	docker-compose run emoji-export

## docker-import:		Imports emoji from the ./import/ directory to your slack team using docker
docker-import: needs-envs
	docker-compose run emoji-import

## download:		Exports emoji from your slack team to the ./export/ directory
download: needs-envs
	pipenv run python main.py export ./export/

## upload:		Imports emoji from the ./import/ directory to your slack team
upload: needs-envs
	pipenv run python main.py import ./import/

stats: needs-envs
	pipenv run python main.py stats

## help:			Prints make target help information from comments in makefile.
help: Makefile
	@sed -n 's/^##//p' $< | sort

needs-envs:
ifndef SLACK_TOKEN
	$(error SLACK_TOKEN is not set)
endif
ifndef SLACK_COOKIE
	$(error SLACK_COOKIE is not set)
endif
ifndef SLACK_TEAM
	$(error SLACK_TEAM is not set)
endif