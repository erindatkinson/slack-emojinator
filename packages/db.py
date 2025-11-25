"""manage DB"""
import os
from datetime import datetime as dt
import peewee
from playhouse.sqlite_ext import SqliteDatabase
from playhouse.migrate import SqliteMigrator, migrate
from .log import get_logger


def get_database()->str:
    """get the database path"""
    return os.getenv("SLACK_EMOJI_DB", './default.sqlite3')


database = SqliteDatabase(get_database())
migrator = SqliteMigrator(database)
# migrations should be in the form of `migrator.add/drop_column(....),`
migrations = []


class Base(peewee.Model):
    """base class for models"""
    created_at = peewee.DateTimeField(default=dt.now)

    class Meta:
        """internal meta class for connecting to the db"""
        database = database


class Download(Base):
    """table to manage downloads"""
    team = peewee.CharField()
    emoji_name = peewee.CharField()


def create_tables():
    """Runs the table create"""
    with database:
        database.create_tables([Download])


def run_migrations():
    """run all the migrations"""
    list(map(migrate, migrations))

def __get_downloaded_emoji(team_name:str):
    """pull the emojis from the database"""
    return Download.select(Download.emoji_name).where(Download.team == team_name).execute()

def connect():
    """connect to the database"""
    database.connect()

def filter_downloaded_emoji(team_name:str, emojis:list)->list:
    """filter out emojis in the list that were already downloaded"""
    downloaded = [emoji.emoji_name for emoji in __get_downloaded_emoji(team_name)]
    return list(filter(lambda emoji: emoji["name"] not in downloaded, emojis))

def mark_emoji_downloaded(team_name:str, emoji: str):
    """mark that an emoji was downloaded"""
    with database.atomic():
        Download.create(team=team_name, emoji_name=emoji)
