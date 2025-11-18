"""Module to handle slack calling functions"""

import asyncio

import os
import time
import aiohttp

from . import log, session


class SlackExportException(Exception):
    """Class for slack export errors"""

class SlackImportException(Exception):
    """Class for slack import errors"""

class NoEmojiException(Exception):
    """Class for empty response errors"""


def post_message(_session:session.Session, channel:str, markdown:str):
    """posts message to slack"""
    message_url = "https://slack.com/api/chat.postMessage"

    params = {
        "token": _session.api_token,
        "channel": channel,
        "as_user": True,
        "markdown_text": markdown + "\n (This was sent via API)"
    }

    return _session.post(message_url, data=params, verify=False)


def upload_emoji(
        _session: session.Session, emoji_name: str, filename: str, logger=log.get_logger()):
    """uploads the emoji data to slack"""

    data = {"mode": "data", "name": emoji_name, "token": _session.api_token}

    while True:
        with open(filename, "rb") as f:
            files = {"image": f}
            resp = _session.post(
                _session.url_add,
                data=data,
                files=files,
                allow_redirects=False,
                verify=False,
            )

            if resp.status_code == 429:
                wait = int(resp.headers.get("retry-after", 1))
                logger.info(f"429 Too Many Requests!, sleeping for {wait} seconds")
                time.sleep(wait)
                continue

        resp.raise_for_status()

        # Slack returns 200 OK even if upload fails, so check for status.
        response_json = resp.json()
        if not response_json["ok"]:
            logger.error(f"Error with uploading {emoji_name}: {response_json}")

        break


def get_current_emoji_list(_session: session.Session, logger=log.get_logger()):
    """List currently uploaded emoji to filter on"""
    page = 1
    result = []
    while True:
        data = {"query": "", "page": page, "count": 1000, "token": _session.api_token}
        logger.debug("Getting emoji list page", url=_session.url_list, page=page)
        resp = _session.post(_session.url_list, data=data, verify=False)
        resp.raise_for_status()
        response_json = resp.json()
        if response_json.get("ok") is False:
            raise SlackImportException(response_json["error"])

        result.extend(map(lambda e: e, response_json["emoji"]))
        if page >= response_json["paging"]["pages"]:
            break

        page = page + 1
    return result

def concurrent_http_get(max_concurrent: int, _session: aiohttp.ClientSession):
    """get emoji data"""
    semaphore = asyncio.Semaphore(max_concurrent)

    async def http_get(emoji: dict):
        nonlocal semaphore
        async with semaphore:
            response = await _session.get(emoji['url'], ssl=False)
            body = await response.content.read()
            await response.wait_for_close()
        return emoji, body

    return http_get


def save_to_file(response: bytes, emoji: dict, directory: str, logger=log.get_logger()):
    """save raw data to file"""
    logger.info("Downloaded %s from %s", emoji['name'].ljust(20), emoji['url'])
    out_fn = f"{emoji['name']}.{str(emoji['url'].rsplit(".", maxsplit=1)[-1])}"
    with open(os.path.join(directory, out_fn), "wb") as out:
        out.write(response)


async def export_loop(
    team_name: str, cookie: str, token: str, directory: str, concurrency: int = 1,
    logger=log.get_logger()
):
    """the async loop for exporting"""

    if not os.path.exists(directory):
        os.makedirs(directory)

    base_session = session.new_session(cookie, team_name, token)
    async with base_session.asyncer() as _session:
        emojis =  get_current_emoji_list(base_session)

        if len(emojis) == 0:
            raise NoEmojiException("Failed to find any custom emoji")

        function_http_get = concurrent_http_get(concurrency, _session)

        for future in asyncio.as_completed(
            [function_http_get(emoji) for emoji in emojis]
        ):
            emoji, data = await future
            save_to_file(data, emoji, directory)

        logger.info(
            "Exported %s custom emoji to directory '%s'", len(emojis), {directory}
        )
