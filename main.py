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

# pylint: disable=import-error
from packages import slack, utils, log
# pylint: enable=import-error


def upload(filepath, cookie: str = "", team: str = "", token: str = ""):
    """Upload a directory of files to the slack team"""
    cookie, team_name, token, _ = utils.arg_envs(cookie, team, token)
    urllib3.disable_warnings(urllib3.exceptions.InsecureRequestWarning)
    logger = log.get_logger()
    session = utils.new_session(cookie, team_name, token)
    existing_emojis = slack.get_current_emoji_list(session)
    uploaded = 0
    skipped = 0

    def process_file(filename):
        nonlocal skipped
        nonlocal uploaded
        logger.info(f"Processing {filename}.")
        emoji_name = f"{os.path.splitext(os.path.basename(filename))[0]}"
        if emoji_name in existing_emojis:
            logger.debug(f"Skipping {emoji_name}. Emoji already exists")
            skipped += 1
        else:
            slack.upload_emoji(session, emoji_name, filename)
            logger.info(f"{filename} upload complete.")
            uploaded += 1

    for slackmoji_file in [filepath]:

        if os.path.isdir(slackmoji_file):
            for file in os.listdir(slackmoji_file):
                filename = os.path.join(slackmoji_file, file)
                process_file(filename)
        else:
            process_file(slackmoji_file)
    logger.info(f"\nUploaded {uploaded} emojis. ({skipped} already existed)")


def stats(cookie: str = "", team: str = "", token: str = ""):
    """getting statistics"""
    cookie, team_name, token, _ = utils.arg_envs(cookie, team, token)
    urllib3.disable_warnings(urllib3.exceptions.InsecureRequestWarning)
    session = utils.new_session(cookie, team_name, token)
    existing_emojis = slack.get_current_emoji_list(session, all_data=True)
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


def export(
    cookie: str = "",
    team: str = "",
    token: str = "",
    export_dir: str = "./export",
    concurrency: int = 1,
):
    """handle the exporting of all emoji from a slack instance"""
    cookie, team_name, token, concurrency = utils.arg_envs(
        cookie, team, token, concurrency=concurrency
    )
    os.makedirs(export_dir, exist_ok=True)
    loop = asyncio.new_event_loop()
    loop.run_until_complete(
        slack.export_loop(
            team_name=team_name, cookie=cookie, token=token, directory=export_dir
        )
    )

def release_notes(
        cookie:str='',
        token:str='',
        team:str='',
        end:str='',
        start:str=str(dt.now() - rdt.relativedelta(days=14))):
    """print the release notes"""

    cookie, team_name, token, _ = utils.arg_envs(cookie, team, token)
    urllib3.disable_warnings(urllib3.exceptions.InsecureRequestWarning)
    session = utils.new_session(cookie, team_name, token)

    _start = dtparse.parse(start)
    if end != '':
        _end = dtparse.parse(end)
    else:
        _end = dt.now()
        
    existing_emojis = slack.get_current_emoji_list(session, all_data=True)

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
    Fire({"export": export, "import": upload, "stats": stats, "release-notes": release_notes})
