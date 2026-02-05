# Slack Emojinator

⚠️ This directory has been kept for posterity, and may not work exactly as described.

![image of heart eyes smiley emoji as a pillow in a field of grass](.github/assets/emoji-4869395_1280.jpg)
Image by [David Bawm](https://pixabay.com/users/david_miram-11502595)

Bulk import and export emoji into Slack

## Setup and Prerequisites

You'll need Python and `pip` to get started. I recommend using [pipenv](https://docs.pipenv.org/).

Prepare a directory that contains an image for each emoji you want to create. Remember to respect Slack's specifications for valid emoji images: no greater than 128px in width or height, no greater than 64K in image size. The base filename of each image file should be the name of the emoji (the bit you'll type inside `:` to display it).

Clone the project and install its prereqs:

```bash
git clone https://github.com/erindatkinson/slack-emojinator.git
cd slack-emojinator
pipenv install
```

You'll need to provide your team name (the bit before ".slack.com" in your admin URL) api token and your session cookie (grab them from your browser). Copy `.env.example`, fill them in, and source it.

### To grab your Slack session cookie and api token

* Open the dev console of your browser and go to the Network tab
* There will be at least 3 calls to /info?, in one of them there will be a payload with a token with custom emoji permissions
* copy that token, and copy the request "as curl" into a text editor
* pull the cookie string from the `-b` flag.

```sh
export SLACK_TEAM=example-team-subdomain
export SLACK_COOKIE='b=<data>; shown_ssb_redirect_page=1; ...; PageCount=99'
export SLACK_TOKEN=xoxc-<numbers>
export SLACK_CONCURRENCY=1
export SLACK_RELEASE_CHANNEL=C0123456789
```

Now you're ready to go!

## Available commands

If you want to know all the commands you can run, just run

```bash
make help
```

```text
 all:                   Sets up pipenv requirements both locally and for docker
 build-docker:          Builds docker image with Pipenv requirements
 docker-export:         Exports emoji from your slack team to the ./export/ directory using docker
 docker-import:         Imports emoji from the ./import/ directory to your slack team using docker
 download:              Exports emoji from your slack team to the ./export/ directory
 help:                  Prints make target help information from comments in makefile.
 lint:                  Runs pylint on the directory
 setup:                 Installs Pipenv requirements locally
 upload:                Imports emoji from the ./import/ directory to your slack team
```

💜
