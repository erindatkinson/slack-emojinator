# Slack Emojinator

![image of heart eyes smiley emoji as a pillow in a field of grass](.github/assets/emoji-4869395_1280.jpg)
Image by [David Bawm](https://pixabay.com/users/david_miram-11502595)

Bulk import and export emoji into Slack

## Setup and Prerequisites

1. Download the correct architecture binary for your computer from the [releases](https://github.com/erindatkinson/slack-emojinator/releases) page.
1. Copy the [.config.yaml.example](.config.yaml.example) to `./.config.yaml` or `$HOME/.emojinator/.config.yaml`
    1. Log into your slack team in your browser if you haven't recently
    1. Configure the yaml appropriately `browser: (chrome, firefox, etc)`
        1. For the profile:
            1. in firefox, `about:profiles` should show you the name, but it's usually `default-release`
            1. in chrome, the `manage profile` page should show a Name field, copy from that and it should be correct.
        1. for the channel, in the slack channel details page is the Channel ID (in the form of `CXXXXXXXXXXX`)
        1. for the subdomain, use the subdomain of your slack team.
    1. if you don't want to store these on your filesystem, you can also put them under the following env vars, or as command line flags:
        1. SLACK_SUBDOMAIN
        1. SLACK_BROWSER
        1. SLACK_PROFILE
        1. SLACK_CHANNEL

## Importing and Exporting emoji

The binary should respond to both the `import` and `export` command.

### Import

Prepare a directory (`./import` is the default) that contains an image for each emoji you want to create. Remember to respect Slack's specifications for valid emoji images: no greater than 128px in width or height, no greater than 64K in image size. The base filename of each image file should be the name of the emoji (the bit you'll type inside `:` to display it).

Run `./slack-emojinator import`, any files with names that already exist in your slack team will be skipped.

### Export

Run `./slack-emojinator export` and the binary should run through all existing emojis and download any emoji that isn't already in your export directory (`./export` is the default).

💜
