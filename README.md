# Slack Emojinator

![image of heart eyes smiley emoji as a pillow in a field of grass](.github/assets/emoji-4869395_1280.jpg)
Image by [David Bawm](https://pixabay.com/users/david_miram-11502595)

Bulk import and export emoji into Slack

## Importing Emoji

You'll need Python and `pip` to get started. I recommend using [pipenv](https://docs.pipenv.org/).

Prepare a directory that contains an image for each emoji you want to create. Remember to respect Slack's specifications for valid emoji images: no greater than 128px in width or height, no greater than 64K in image size. The base filename of each image file should be the name of the emoji (the bit you'll type inside `:` to display it).

Clone the project and install its prereqs:

```bash
git clone https://github.com/erindatkinson/slack-emojinator.git
cd slack-emojinator
pipenv install
```

You'll need to provide your team name (the bit before ".slack.com" in your admin URL) api token and your session cookie (grab them from your browser). Copy `.env.example`, fill them in, and source it.

To grab your Slack session cookie and api token:

* [Open your browser's dev tools](http://webmasters.stackexchange.com/a/77337) and copy the value of `document.cookie`.
* Go to the Network tab.
* Re-load your workspace's `https://{teamname}.slack.com/customize/emoji` page.
* Find the various calls to `info`.
  * In the payload of one of them will be the token, copy the value of the token and add to your `.env` file.
  * Scroll to `Request-Headers`, copy the value of "Cookie,"and add to your `.env` file.
    * You can also "copy request as curl and copy the string after `-b`

```bash
cp .env.example .env
${EDITOR} .env
source .env
```

Now you're ready to go!

```bash
make upload
```

:sparkles:

## Exporting Emoji

To export emoji

```bash
source .env
make download
```

## Docker

If you'd rather run this through docker, you can run everything the same as before, ensuring you sourced a filled .env file:

### Docker: Import

Put emoji in `./import`

```bash
source .env
make docker-import
```

### Docker: Export

Put emoji in `./export`

```bash
source .env
make docker-export
```

## Available commands

If you want to know all the commands you can run, just run

```bash
make help
```

💜
