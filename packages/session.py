"""module for handling sessions"""
import os
import urllib3
import requests
import aiohttp


URL_CUSTOMIZE = "https://{team_name}.slack.com/customize/emoji"
URL_ADD = "https://{team_name}.slack.com/api/emoji.add"
URL_LIST = "https://{team_name}.slack.com/api/emoji.adminList"
POST_MESSAGE = "https://slack.com/api/chat.postMessage"



class Session(requests.Session):
    """wrapper class around requests.Session for additional
    attributes"""

    def __init__(self, team_name, token, cookie, concurrency):
        super().__init__()
        self.headers = {"Cookie": cookie}
        self.url_customize = URL_CUSTOMIZE.format(team_name=team_name)
        self.url_add = URL_ADD.format(team_name=team_name)
        self.url_list = URL_LIST.format(team_name=team_name)
        self.api_token = token
        self.concurrency = concurrency

    def asyncer(self)-> aiohttp.ClientSession:
        """give an async session from this session"""
        return aiohttp.ClientSession(headers=self.headers, trust_env=True)


def new_session() -> Session:
    """Set up session object for making requests with slack cookie/token"""
    cookie = os.getenv("SLACK_COOKIE")
    team_name = os.getenv("SLACK_TEAM")
    token = os.getenv("SLACK_TOKEN")
    concurrency = int(os.getenv("SLACK_CONCURRENCY", '1'))

    assert cookie, "Either SLACK_COOKIE env var, or --cookie param must be set"
    assert team_name, "Either SLACK_TEAM env var, or --team param must be set"
    assert token, "Either SLACK_TOKEN env var, or --token param must be set"

    session = Session(team_name, token, cookie, concurrency)
    urllib3.disable_warnings(urllib3.exceptions.InsecureRequestWarning)

    return session
