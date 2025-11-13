"""Module to handle sessions"""

import os

import requests
import aiohttp


BASE_URL = "https://{team_name}.slack.com"
EMOJI_ENDPOINT = "/customize/emoji"
EMOJI_API = "/api/emoji.adminList"

URL_CUSTOMIZE = "https://{team_name}.slack.com/customize/emoji"
URL_ADD = "https://{team_name}.slack.com/api/emoji.add"
URL_LIST = "https://{team_name}.slack.com/api/emoji.adminList"


class Session(requests.Session):
    """wrapper class around requests.Session for additional
    attributes"""

    def __init__(self, team_name, token):
        super().__init__()
        self.url_customize = URL_CUSTOMIZE.format(team_name=team_name)
        self.url_add = URL_ADD.format(team_name=team_name)
        self.url_list = URL_LIST.format(team_name=team_name)
        self.api_token = token


def new_session(cookie: str, team_name: str, token: str) -> Session:
    """Set up session object for making requests with slack cookie/token"""
    assert cookie, "Cookie required"
    assert team_name, "Team name required"
    session = Session(team_name, token)
    session.headers = {"Cookie": cookie}

    return session


def async_session(auth_cookie) -> aiohttp.ClientSession:
    """create a session object for the async runs"""
    return aiohttp.ClientSession(headers={"Cookie": auth_cookie})


def arg_envs(cookie: str, team_name: str, token: str, concurrency: int = 1) -> tuple:
    """Pull args from env if needed, and assert requirements"""
    _cookie = os.getenv("SLACK_COOKIE", cookie)
    _team_name = os.getenv("SLACK_TEAM", team_name)
    _token = os.getenv("SLACK_TOKEN", token)
    _concurrency = int(os.getenv("SLACK_CONCURRENCY", concurrency))

    assert _cookie, "Either SLACK_COOKIE env var, or --cookie param must be set"
    assert _team_name, "Either SLACK_TEAM env var, or --team param must be set"
    assert _token, "Either SLACK_TOKEN env var, or --token param must be set"
    assert (
        _concurrency
    ), "Either SLACK_CONCURRENCY env var, or --concurrency param must be set"

    return (_cookie, _team_name, _token, _concurrency)


def cookie_split(cookie: str = os.getenv("SLACK_COOKIE", "")):
    """method to split a cookie string int a list of kv pairs"""
    if cookie == "":
        return []
    return list(list(map(lambda x: x.split("="), cookie.split(";"))))
