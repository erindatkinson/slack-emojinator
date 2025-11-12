#!/usr/bin/env python
"""An application to bulk upload emoji to slack."""

# https://github.com/smashwilson/slack-emojinator

from __future__ import print_function

import os
import os.path
from tabulate import tabulate
import urllib3

from fire import Fire
import numpy as np

from packages import slack, utils

def upload(filepath, cookie:str='', team:str='', token:str='' ):
    """Upload a directory of files to the slack team"""
    cookie, team_name, token = utils.arg_envs(cookie, team, token)
    urllib3.disable_warnings(urllib3.exceptions.InsecureRequestWarning)
    session = utils.new_session(cookie, team_name, token)
    existing_emojis = slack.get_current_emoji_list(session)
    uploaded = 0
    skipped = 0

    def process_file(filename):
        nonlocal skipped
        nonlocal uploaded
        print(f"Processing {filename}.")
        emoji_name = f"{os.path.splitext(os.path.basename(filename))[0]}"
        if emoji_name in existing_emojis:
            print(f"Skipping {emoji_name}. Emoji already exists")
            skipped += 1
        else:
            slack.upload_emoji(session, emoji_name, filename)
            print(f"{filename} upload complete.")
            uploaded += 1

    for slackmoji_file in [filepath]:

        if os.path.isdir(slackmoji_file):
            for file in os.listdir(slackmoji_file):
                print()
                filename = os.path.join(slackmoji_file, file)
                process_file(filename)
        else:
            process_file(slackmoji_file)
    print(f"\nUploaded {uploaded} emojis. ({skipped} already existed)")

def stats(cookie:str='', team:str='', token:str=''):
    """getting statistics"""
    cookie, team_name, token = utils.arg_envs(cookie, team, token)
    urllib3.disable_warnings(urllib3.exceptions.InsecureRequestWarning)
    session = utils.new_session(cookie, team_name, token)
    existing_emojis = slack.get_current_emoji_list(session, all=True)
    userstats = {}
    for emoji in existing_emojis:
        try:
            userstats[emoji['user_display_name']] += 1
        except KeyError:
            userstats[emoji['user_display_name']] = 1
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


if __name__=="__main__":
    Fire({
        "upload": upload,
        "stats": stats
    })
