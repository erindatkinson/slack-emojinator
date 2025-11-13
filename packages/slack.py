"""Module to handle slack calling functions"""

import asyncio
import collections
import os
import time
import aiohttp

from . import utils, log


class SlackExportException(Exception):
    """Class for slack export errors"""


class NoEmojiException(Exception):
    """Class for empty response errors"""


Emoji = collections.namedtuple("Emoji", ["url", "name", "extension"])


def upload_emoji(session: utils.Session, emoji_name: str, filename: str):
    """uploads the emoji data to slack"""
    logger = log.get_logger()
    data = {"mode": "data", "name": emoji_name, "token": session.api_token}

    while True:
        with open(filename, "rb") as f:
            files = {"image": f}
            resp = session.post(
                session.url_add,
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


def get_current_emoji_list(session: utils.Session, all_data=False):
    """List currently uploaded emoji to filter on"""
    page = 1
    result = []
    while True:
        data = {"query": "", "page": page, "count": 1000, "token": session.api_token}
        resp = session.post(session.url_list, data=data, verify=False)
        resp.raise_for_status()
        response_json = resp.json()

        if not all_data:
            result.extend(map(lambda e: e["name"], response_json["emoji"]))
            if page >= response_json["paging"]["pages"]:
                break
        else:
            result.extend(response_json["emoji"])
            if page >= response_json["paging"]["pages"]:
                break

        page = page + 1
    return result


async def _determine_all_emoji_urls(
    session: aiohttp.ClientSession, base_url: str, token: str
):
    """pull all urls to download all emoji"""
    logger = log.get_logger()
    page = 1
    total_pages = None

    entries = []

    while total_pages is None or page <= total_pages:

        data = {"token": token, "page": page, "count": 1000}

        response = await session.post(base_url + utils.EMOJI_API, data=data, ssl=False)

        logger.info("loaded %s (page %d)", response.real_url, page)

        if response.status == 429:
            retry = int(response.headers.get("RETRY-AFTER", 60))
            time.sleep(retry)
            response = await session.post(base_url + utils.EMOJI_API, data=data, ssl=False)

            logger.info("Rate-Limited: loaded %s (page %d)", response.real_url, page)
        elif response.status != 200:
            real_url = response.request_info.real_url
            raise SlackExportException(
                f"Failed to load emoji from {real_url} (status {response.status})"
            )

        json = await response.json()

        for entry in json["emoji"]:
            url = str(entry["url"])
            name = str(entry["name"])
            extension = str(url.rsplit(".", maxsplit=1)[-1])
            entries.append(Emoji(url, name, extension))

        if total_pages is None:
            total_pages = int(json["paging"]["pages"])

        page += 1

    return entries


def concurrent_http_get(max_concurrent: int, session: aiohttp.ClientSession):
    """get emoji data"""
    semaphore = asyncio.Semaphore(max_concurrent)

    async def http_get(emoji: Emoji):
        nonlocal semaphore
        async with semaphore:
            response = await session.get(emoji.url, ssl=False)
            body = await response.content.read()
            await response.wait_for_close()
        return emoji, body

    return http_get


def save_to_file(response: bytes, emoji: Emoji, directory: str):
    """save raw data to file"""
    logger = log.get_logger()
    logger.info("Downloaded %s from %s", emoji.name.ljust(20), emoji.url)
    with open(os.path.join(directory, f"{emoji.name}.{emoji.extension}"), "wb") as out:
        out.write(response)


async def export_loop(
    team_name: str, cookie: str, token: str, directory: str, concurrency: int = 1
):
    """the async loop for exporting"""
    logger = log.get_logger()

    if not os.path.exists(directory):
        os.makedirs(directory)

    base_url = utils.BASE_URL.format(team_name=team_name)

    async with utils.async_session(cookie) as session:
        emojis = await _determine_all_emoji_urls(session, base_url, token)

        if len(emojis) == 0:
            raise NoEmojiException("Failed to find any custom emoji")

        function_http_get = concurrent_http_get(concurrency, session)

        for future in asyncio.as_completed(
            [function_http_get(emoji) for emoji in emojis]
        ):
            emoji, data = await future
            save_to_file(data, emoji, directory)

        logger.info(
            "Exported %s custom emoji to directory '%s'", len(emojis), {directory}
        )
