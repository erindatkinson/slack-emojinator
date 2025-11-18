"""Module to handle sessions"""

import os

from jinja2 import FileSystemLoader, Environment


def preprocess_slackmoji(slackmoji):
    """Process the slackmoji to get fn and name"""
    if os.path.isdir(slackmoji):
        for file in os.listdir(slackmoji):
            filename = os.path.join(slackmoji, file)
    else:
        filename = slackmoji

    return [filename]


def arg_envs() -> tuple:
    """Pull args from env if needed, and assert requirements"""
    cookie = os.getenv("SLACK_COOKIE")
    team_name = os.getenv("SLACK_TEAM")
    token = os.getenv("SLACK_TOKEN")
    concurrency = int(os.getenv("SLACK_CONCURRENCY", '1'))

    assert cookie, "Either SLACK_COOKIE env var, or --cookie param must be set"
    assert team_name, "Either SLACK_TEAM env var, or --team param must be set"
    assert token, "Either SLACK_TOKEN env var, or --token param must be set"
    assert (
        concurrency
    ), "Either SLACK_CONCURRENCY env var, or --concurrency param must be set"

    return (cookie, team_name, token, concurrency)


def load_templates(tpl_dir:str)->Environment:
    """Load templates from directory"""
    return Environment(loader=FileSystemLoader([tpl_dir, './templates']))


def filter_emojis_to_span(start, end, current)->list:
    """filters current emoji list between times"""

    return list(
                filter(
                    lambda x:
                        x['created'] <= end.timestamp() and x['created'] > start.timestamp(),
                    sorted(
                        current,
                        key=lambda x: x['name'])))

def format_emojis_into_string_list(emojis:list)->list:
    """create list of formatted strings from emojis"""
    return list(map(lambda x: f"* :{x['name']}: | `{x['name']}`", emojis))

def build_user_ranks(emojis:list)->list:
    "build a sorted list of kv pairs (name, upload_count)"
    ranks = {}
    for item in emojis:
        try:
            ranks[item['user_display_name']] += 1
        except KeyError:
            ranks[item['user_display_name']] = 1
    return sorted(ranks.items(), key=lambda x: x[1], reverse=True)
