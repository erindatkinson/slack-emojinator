"""Module to handle slack calling functions"""

import asyncio

import os
import time
from requests import Response
from requests.exceptions import HTTPError

from . import log, session, utils


class SlackExportException(Exception):
    """Class for slack export errors"""

class SlackImportException(Exception):
    """Class for slack import errors"""

class NoEmojiException(Exception):
    """Class for empty response errors"""

class Slack:
    """Slack client wrapper"""
    def __init__(
            self, _session:session.Session=None, logger=None):
        self.session = _session
        self.api_token = _session.api_token
        self.concurrency = _session.concurrency

        if logger is not None:
            self.log = logger
        else:
            self.log = log.get_logger()


    def post_message(
            self, message:str, channel:str, thread_ts:str|None=None, verify:bool=False) -> Response:
        """Post a markdown message to a slack channel"""

        data = {
            "token": self.api_token,
            "channel": channel,
            "as_user": True,
            "markdown_text": message + "\n (This was sent via API)"
        }

        if thread_ts is not None:
            data.update({"thread_ts": thread_ts})

        return self.session.post(session.POST_MESSAGE, data=data, verify=verify)


    def list_emoji(self):
        """List currently uploaded emoji to filter on"""
        page = 1
        result = []
        while True:
            data = {"query": "", "page": page, "count": 1000, "token": self.api_token}
            self.log.debug("Getting emoji list page", url=self.session.url_list, page=page)
            resp = self.session.post(self.session.url_list, data=data, verify=False)
            resp.raise_for_status()
            response_json = resp.json()
            if response_json.get("ok") is False:
                raise SlackImportException(response_json["error"])

            result.extend(map(lambda e: e, response_json["emoji"]))
            if page >= response_json["paging"]["pages"]:
                break

            page = page + 1
        return result


    def import_emoji(self, filepath:str):
        """Run loop for importing all emoji in a filepath"""
        try:
            existing_emojis = self.list_emoji()
        except SlackImportException as sie:
            self.log.error("Unable to get current emojis", error=sie)
            return

        for filename in utils.preprocess_slackmoji(filepath):
            emoji_name = f"{os.path.splitext(os.path.basename(filename))[0]}"
            self.log.info(f"Processing {filename}.")

            if emoji_name in existing_emojis:
                self.log.debug(f"Skipping {emoji_name}. Emoji already exists")
                continue
            else:
                try:
                    self.__upload_emoji(emoji_name, filename)
                    self.log.info(f"{filename} upload complete.")
                except HTTPError as he:
                    self.log.error("Bad response status when uploading", error=he)


    async def export_emoji(self, directory:str):
        """Export emoji asynchronously"""
        if not os.path.exists(directory):
            os.makedirs(directory)
        emojis = self.list_emoji()
        if len(emojis) == 0:
            raise NoEmojiException("Failed to find any custom emoji")
        function_http_get = self.__concurrent_http_get()
        for future in asyncio.as_completed([function_http_get(emoji) for emoji in emojis]):
            emoji, data = await future
            self.__save_to_file(data, emoji, directory)

        self.log.info(
            "Exported %s custom emoji to directory '%s'", len(emojis), {directory}
        )


    def __save_to_file(self, response:bytes, emoji:dict, directory:str):
        """Save the raw data to a file"""
        self.log.info("Downloaded %s from %s", emoji['name'].ljust(20), emoji['url'])
        out_fn = f"{emoji['name']}.{str(emoji['url'].rsplit(".", maxsplit=1)[-1])}"
        with open(os.path.join(directory, out_fn), "wb") as out:
            out.write(response)


    def __concurrent_http_get(self):
        """get emoji data"""
        semaphore = asyncio.Semaphore(self.concurrency)

        async def http_get(emoji: dict):
            nonlocal semaphore
            async with semaphore:
                response = await self.session.asyncer().get(emoji['url'], ssl=False)
                body = await response.content.read()
                await response.wait_for_close()
            return emoji, body

        return http_get


    def __upload_emoji(self, emoji_name: str, filename: str):
        """uploads the emoji data to slack"""

        data = {"mode": "data", "name": emoji_name, "token": self.api_token}

        while True:
            with open(filename, "rb") as f:
                files = {"image": f}
                resp = self.session.post(
                    self.session.url_add,
                    data=data,
                    files=files,
                    allow_redirects=False,
                    verify=False,
                )

                if resp.status_code == 429:
                    wait = int(resp.headers.get("retry-after", 1))
                    self.log.info(f"429 Too Many Requests!, sleeping for {wait} seconds")
                    time.sleep(wait)
                    continue

            resp.raise_for_status()

            # Slack returns 200 OK even if upload fails, so check for status.
            response_json = resp.json()
            if not response_json["ok"]:
                self.log.error(f"Error with uploading {emoji_name}: {response_json}")

            break
