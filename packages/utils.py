"""Module to handle sessions"""

import os
from datetime import datetime as dt
import dateutil.parser as dtparse
from jinja2 import FileSystemLoader, Environment
import numpy as np
from numpy.typing import NDArray
from .db import create_tables, run_migrations, get_database

def setup_duration_span(start, end)->tuple:
    """set the duration for running stats gathering"""
    _start = dtparse.parse(start)
    if end != '':
        _end = dtparse.parse(end)
    else:
        _end = dt.now()
    return (_start, _end)


def process_stats(emojis)-> tuple[dict, NDArray]:
    """Build userstats"""
    userstats = {}
    for emoji in emojis:
        try:
            userstats[emoji["user_display_name"]] += 1
        except KeyError:
            userstats[emoji["user_display_name"]] = 1

    df = np.array(list(userstats.values()))
    return(userstats, df)


def preprocess_slackmoji(slackmoji):
    """Process the slackmoji to get fn and name"""
    if os.path.isdir(slackmoji):
        files = []
        for file in os.listdir(slackmoji):
            files.append(os.path.join(slackmoji, file))
        return files
    else:
        return [slackmoji]


def load_templates(tpl_dir:str)->Environment:
    """Load templates from directory"""
    return Environment(loader=FileSystemLoader([tpl_dir, './templates']))


def filter_emojis_to_span(start:dt, end:dt, emojis:list)->list:
    """filters current emoji list between times"""

    return list(
                filter(
                    lambda x:
                        x['created'] <= end.timestamp() and x['created'] > start.timestamp(),
                    sorted(
                        emojis,
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

def setup(migrate:bool=False):
    """various setup tasks"""
    if not os.path.isfile(get_database()):
        create_tables()
        run_migrations()
    else:
        if migrate:
            run_migrations()
