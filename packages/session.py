"""module for handling sessions"""
import urllib3
import requests
import aiohttp


URL_CUSTOMIZE = "https://{team_name}.slack.com/customize/emoji"
URL_ADD = "https://{team_name}.slack.com/api/emoji.add"
URL_LIST = "https://{team_name}.slack.com/api/emoji.adminList"



class Session(requests.Session):
    """wrapper class around requests.Session for additional
    attributes"""

    def __init__(self, team_name, token, cookie):
        super().__init__()
        self.headers = {"Cookie": cookie}
        self.url_customize = URL_CUSTOMIZE.format(team_name=team_name)
        self.url_add = URL_ADD.format(team_name=team_name)
        self.url_list = URL_LIST.format(team_name=team_name)
        self.api_token = token

    def asyncer(self)-> aiohttp.ClientSession:
        """give an async session from this session"""
        return aiohttp.ClientSession(headers=self.headers, trust_env=True)


def new_session(cookie: str, team_name: str, token: str) -> Session:
    """Set up session object for making requests with slack cookie/token"""
    assert cookie, "Cookie required"
    assert team_name, "Team name required"
    session = Session(team_name, token, cookie)
    urllib3.disable_warnings(urllib3.exceptions.InsecureRequestWarning)

    return session
