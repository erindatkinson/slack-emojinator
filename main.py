#!/usr/bin/env python
"""An application to bulk upload emoji to slack."""

# https://github.com/erindatkinson/slack-emojinator

import asyncio
import os
import os.path
import sys
import urllib3
from tabulate import tabulate
from selenium import webdriver
from selenium.webdriver.support.wait import WebDriverWait

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


def testing(browser: str = "firefox"):
    """Testing selenium option"""
    logger = log.get_logger()
    driver_opts = {
        "firefox": webdriver.Firefox,
        "chrome": webdriver.Chrome,
        "safari": webdriver.Safari,
    }

    try:
        driver = driver_opts[browser]()
        team_name = os.getenv("SLACK_TEAM")
        cookie_string = os.getenv("SLACK_COOKIE")
        assert cookie_string, "SLACK_COOKIE must be set"
        assert team_name, "SLACK_TEAM must be set"

        # have to call it first to load the context before adding a cookie
        driver.get(utils.URL_CUSTOMIZE.format(team_name=team_name))
        cookies = utils.cookie_split(cookie_string)

        for cookie in cookies:
            driver.add_cookie({"name": cookie[0].strip(), "value": cookie[1].strip()})

        # actually get the thing
        driver.get(utils.URL_CUSTOMIZE.format(team_name=team_name))
        driver.implicitly_wait(2)

        wait = WebDriverWait(driver, timeout=5)
        api_token = wait.until(
            lambda driver: driver.execute_script("return boot_data.api_token;")
        )

        print(api_token)
    except KeyError:
        logger.error(
            f"{browser} is unsupported, please use one of [{', '.join(driver_opts.keys())}]"
        )
        sys.exit(1)
    finally:
        if driver:  # type: ignore
            driver.close()


if __name__ == "__main__":
    Fire({"export": export, "import": upload, "stats": stats, "test": testing})
