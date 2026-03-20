# Emoji Archiver

![image of heart eyes smiley emoji as a pillow in a field of grass](.github/assets/emoji-4869395_1280.jpg)
Image by [David Bawm](https://pixabay.com/users/david_miram-11502595)

Bulk import and export emoji into Slack

## Setup and Prerequisites

1. Download the correct architecture binary for your computer from the [releases](https://github.com/erindatkinson/emoji-archiver/releases) page.
1. Copy the [.config.yaml.example](.config.yaml.example) to `./.config.yaml` or `$HOME/.emojinator/.config.yaml`
    1. Log into your slack team in your browser if you haven't recently
    1. Configure the yaml appropriately `browser: (chrome, firefox, etc)`
        1. For the profile:
            1. in firefox, `about:profiles` should show you the name, but it's usually `default-release`
            1. in chrome, the `Customize Profile` or `Edit` profile page should show a "Name your Chrome profile" field, copy from that and it should be correct, it may be different from what's shown in the popout.
        1. for the channel, in the slack channel details page is the Channel ID (in the form of `CXXXXXXXXXXX`)
        1. for the subdomain, use the subdomain of your slack team.
    1. if you don't want to store these on your filesystem, you can also put them under the following env vars, or as command line flags:
        1. SLACK_SUBDOMAIN
        1. SLACK_BROWSER
        1. SLACK_PROFILE
        1. SLACK_CHANNEL

### Listing available profile/browser combinations

Running `./emoji-archiver list-profiles` with a subdomain configured either in config.yaml, or with the `-s` flag will output something like the following:

```shell
╭─────────┬─────────────────╮
│ BROWSER │ PROFILE         │
├─────────┼─────────────────┤
│ chrome  │ Your Name       │
│ firefox │ default-release │
╰─────────┴─────────────────╯
```

where each row is a browser/profile combination that has a cookie stored for that subdomain.

## Importing and Exporting emoji

The binary should respond to both the `import` and `export` command.

### Import

Prepare a directory (`./emojis/<subdomain>` is the default) that contains an image for each emoji you want to create. Remember to respect Slack's specifications for valid emoji images: no greater than 128px in width or height, no greater than 64K in image size. The base filename of each image file should be the name of the emoji (the bit you'll type inside `:` to display it).

Run `./emoji-archiver import`, any files with names that already exist in your slack team will be skipped.

### Export

Run `./emoji-archiver export` and the binary should run through all existing emojis and download any emoji that isn't already in your export directory (`./emojis/<subdomain>` is the default).

## Generating Docs Markdown

Run `./emoji-archiver docs` and the binary should generate an index file and pages of 100 emojis.
The emojis will be generated based off of `./emojis/<subdomain>` by default and be populated in `./docs/<subdomain>` by default.

## Posting "Emoji Release Notes" for a Slack team

Running `./emoji-archiver release-notes` will post a ranking of emoji uploaders, and a sorted list of new emojis to the configured .slack.channel option in the .config.yaml

💜
