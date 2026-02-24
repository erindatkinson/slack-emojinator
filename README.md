# Slack Emojinator

![image of heart eyes smiley emoji as a pillow in a field of grass](.github/assets/emoji-4869395_1280.jpg)
Image by [David Bawm](https://pixabay.com/users/david_miram-11502595)

Bulk import and export emoji into Slack

## Setup and Prerequisites

1. Download the correct architecture binary for your computer from the [releases](https://github.com/erindatkinson/slack-emojinator/releases) page.
1. Set the required env vars described in [`.env.example`](.env.example)
    1. Steps to acquire the info are detailed below

### To grab your Slack session cookie and api token

* Navigate to the Custom Emoji Page of your slack team
* Open the dev console of your browser and go to the Network tab (you may have to refresh)
* There will be at least 3 calls to /info?, in one of them there will be a payload with a token with custom emoji permissions
* copy that token, and copy the request "as curl" into a text editor
* pull the cookie string from the `-b` flag if Chrome, or the data after `-H "Cookie` if Firefox.

```sh
export SLACK_TEAM=example-team-subdomain
export SLACK_COOKIE='b=<data>; shown_ssb_redirect_page=1; ...; PageCount=99'
export SLACK_TOKEN=xoxc-<numbers>
export SLACK_CONCURRENCY=1
export SLACK_RELEASE_CHANNEL=C0123456789
```

## Importing and Exporting emoji

The binary should respond to both the `import` and `export` command.

### Import

Prepare a directory (`./import` is the default) that contains an image for each emoji you want to create. Remember to respect Slack's specifications for valid emoji images: no greater than 128px in width or height, no greater than 64K in image size. The base filename of each image file should be the name of the emoji (the bit you'll type inside `:` to display it).

Run `./slack-emojinator import`, any files with names that already exist in your slack team will be skipped.

### Export

Run `./slack-emojinator export` and the binary should run through all existing emojis and download any emoji that isn't already in your export directory (`./export` is the default).

💜
