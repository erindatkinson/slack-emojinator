#!/usr/bin/env python
"""An application to bulk upload emoji to slack."""

# https://github.com/erindatkinson/slack-emojinator

import asyncio
import os
import os.path
from datetime import datetime as dt
import dateutil.relativedelta as rdt
import dateutil.parser as dtparse

import urllib3
from tabulate import tabulate
from fire import Fire
import numpy as np
from requests.exceptions import HTTPError

# pylint: disable=import-error
from packages import slack, utils, log, session
# pylint: enable=import-error


def import_emoji(filepath):
    """Upload a directory of files to a slack team"""
    logger = log.get_logger()
    cookie, team_name, token, _ = utils.arg_envs()
    _session = session.new_session(cookie, team_name, token)
    try:
        existing_emojis = slack.get_current_emoji_list(_session)
    except slack.SlackImportException as sie:
        logger.error("Unable to get current emojis", error=sie)
        return

    logger.debug("")
    for filename in utils.preprocess_slackmoji(filepath):
        emoji_name = f"{os.path.splitext(os.path.basename(filename))[0]}"
        logger.info(f"Processing {filename}.")

        if emoji_name in existing_emojis:
            logger.debug(f"Skipping {emoji_name}. Emoji already exists")
            continue
        else:
            try:
                slack.upload_emoji(_session, emoji_name, filename, logger)
                logger.info(f"{filename} upload complete.")
            except HTTPError as he:
                logger.error("Bad response status when uploading", error=he)


def stats():
    """getting statistics"""
    cookie, team_name, token, _ = utils.arg_envs()
    urllib3.disable_warnings(urllib3.exceptions.InsecureRequestWarning)
    _session = session.new_session(cookie, team_name, token)
    existing_emojis = slack.get_current_emoji_list(_session)
    userstats = {}
    for emoji in existing_emojis:
        try:
            userstats[emoji["user_display_name"]] += 1
        except KeyError:
            userstats[emoji["user_display_name"]] = 1
    print(tabulate(sorted(userstats.items(), key=lambda x: x[1], reverse=True)[:25]))

    df = np.array(list(userstats.values()))
    ninety_nine_q = np.percentile(df, 99)
    print(f"99th Percentile: {ninety_nine_q}")
    ninety_q = np.percentile(df, 90)
    print(f"90th Percentile: {ninety_q}")
    q1 = np.percentile(df, 75)
    print(f"Top Quartile: {q1}")
    q2 = np.percentile(df, 50)
    print(f"Middle Quartile: {q2}")
    q3 = np.percentile(df, 25)
    print(f"Bottom Quartile: {q3}")


def export_emoji(export_dir: str = "./export"):
    """handle the exporting of all emoji from a slack instance"""
    cookie, team_name, token, concurrency = utils.arg_envs()
    os.makedirs(export_dir, exist_ok=True)

    loop = asyncio.new_event_loop()
    loop.run_until_complete(
        slack.export_loop(
            team_name=team_name,
            cookie=cookie,
            token=token,
            directory=export_dir,
            concurrency=concurrency
        )
    )

def release_notes(end:str='', start:str=str(dt.now() - rdt.relativedelta(days=14))):
    """print the release notes"""

    cookie, team_name, token, _ = utils.arg_envs()
    _session = session.new_session(cookie, team_name, token)

    _start = dtparse.parse(start)
    if end != '':
        _end = dtparse.parse(end)
    else:
        _end = dt.now()

    existing_emojis = slack.get_current_emoji_list(_session)

    span = list(
                filter(
                    lambda x:
                        x['created'] <= _end.timestamp() and x['created'] > _start.timestamp(),
                    sorted(
                        existing_emojis,
                        key=lambda x: x['name'])))

    list_items = list(map(lambda x: f"* :{x['name']}: | `{x['name']}`", span))

    ranks = {}
    for item in span:
        try:
            ranks[item['user_display_name']] += 1
        except KeyError:
            ranks[item['user_display_name']] = 1
    sorted_ranks = sorted(ranks.items(), key=lambda x: x[1], reverse=True)


    tpls = utils.load_templates(".")
    rn_tpl = tpls.get_template("release_notes.md.jinja2")
    print(rn_tpl.render(
        start=_start.strftime("%Y-%m-%d"),
        end=_end.strftime("%Y-%m-%d"),
        ranks=tabulate(sorted_ranks),
        emojis=list_items
        )
    )




if __name__ == "__main__":
    Fire({
        "export": export_emoji, 
        "import": import_emoji, 
        "stats": stats, 
        "release-notes": release_notes
        })
