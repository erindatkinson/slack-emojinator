#!/usr/bin/env python
"""An application to bulk upload emoji to slack."""

# https://github.com/erindatkinson/slack-emojinator

import asyncio
import os
import os.path
from datetime import datetime as dt
import dateutil.relativedelta as rdt


from tabulate import tabulate
from fire import Fire
import numpy as np

# pylint: disable=import-error
from packages import slack, utils, log, session
# pylint: enable=import-error


def import_emoji(filepath:str):
    """Import all emoji in the given filepath to the connected slack team"""
    client = slack.Slack(session.new_session(), log.get_logger())
    client.import_emoji(filepath)


def export_emoji(export_dir: str = "./export"):
    """Export all emoji in the connected slack team into the given export_dir"""
    client = slack.Slack(session.new_session(), log.get_logger())
    os.makedirs(export_dir, exist_ok=True)

    loop = asyncio.new_event_loop()
    loop.run_until_complete(
        client.export_emoji(directory=export_dir)
    )


def release_notes(end:str='', start:str=str(dt.now() - rdt.relativedelta(days=14))):
    """Retrieve the emojis uploaded within the span of time and, if within Slack's
    message limits, post formatted list to a slack channel, otherwise, print to
    standard out"""

    # set up the prereq variables
    logger = log.get_logger()
    channel = os.getenv("SLACK_RELEASE_CHANNEL")
    client = slack.Slack(session.new_session(), logger)

    # pull current emoji json details, filter, and build output lists
    _start, _end = utils.setup_duration_span(start, end)
    start_str, end_str = _start.strftime("%Y-%m-%d"), _end.strftime("%Y-%m-%d")
    existing_emojis = client.list_emoji()
    span = utils.filter_emojis_to_span(_start, _end, existing_emojis)
    list_items = utils.format_emojis_into_string_list(span)
    ranks = utils.build_user_ranks(span)

    # load jinja templates
    tpls = utils.load_templates(".")
    rn_tpl = tpls.get_template("release_notes.md.jinja2")
    header_tpl = tpls.get_template("header.md.jinja2")

    # render output
    header = header_tpl.render(start=start_str, end=end_str)
    markdown = rn_tpl.render(ranks=tabulate(ranks), emojis=list_items)

    # post to slack or print to local
    if len(markdown) <= 12_000:
        resp = client.post_message(header, channel)
        try:
            if resp.status_code == 200:
                client.post_message(markdown, channel, resp.json()["ts"])
        except KeyError:
            logger.error("failed posting thread message", json=resp.json())
    else:
        logger.warn("message is over slack's posting limit of 12,000 characters",
                    length=len(markdown))
        print(markdown)


def stats():
    """Calculate some statistics about the slack team, mainly
    top 25 uploaders, and the percentile breakdown for uploads"""
    logger = log.get_logger()
    client = slack.Slack(session.new_session(), logger)
    existing_emojis = client.list_emoji()
    userstats, df = utils.process_stats(existing_emojis)

    tpls = utils.load_templates(".")
    tpls.globals['tabulate'] = tabulate
    stats_tpl = tpls.get_template("stats.md.jinja2")
    print(stats_tpl.render(
        top_25=sorted(userstats.items(), key=lambda x: x[1], reverse=True)[:25],
        pc_99=np.percentile(df, 99),
        pc_90=np.percentile(df, 90),
        pc_75=np.percentile(df, 75),
        pc_50=np.percentile(df, 50),
        pc_25=np.percentile(df, 25),
    ))


if __name__ == "__main__":
    Fire({
        "export": export_emoji, 
        "import": import_emoji, 
        "stats": stats, 
        "release-notes": release_notes,
        "setup": utils.setup
        })
