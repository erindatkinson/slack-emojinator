"""Module to handle sessions"""
import os

import requests
import aiohttp


BASE_URL = 'https://{team_name}.slack.com'
EMOJI_ENDPOINT = '/customize/emoji'
EMOJI_API = '/api/emoji.adminList'

URL_CUSTOMIZE = "https://{team_name}.slack.com/customize/emoji"
URL_ADD = "https://{team_name}.slack.com/api/emoji.add"
URL_LIST = "https://{team_name}.slack.com/api/emoji.adminList"

def new_session(cookie:str, team_name:str, token:str)-> requests.Session:
    """Set up session object for making requests with slack cookie/token"""
    assert cookie, "Cookie required"
    assert team_name, "Team name required"
    session = requests.session()
    session.headers = {'Cookie': cookie}
    session.url_customize = URL_CUSTOMIZE.format(team_name=team_name)
    session.url_add = URL_ADD.format(team_name=team_name)
    session.url_list = URL_LIST.format(team_name=team_name)
    session.api_token = token
    return session


def async_session(auth_cookie) -> aiohttp.ClientSession:
    """create a session object for the async runs"""
    return aiohttp.ClientSession(headers={"Cookie": auth_cookie})


# For export
# def _argparse():
#     parser = argparse.ArgumentParser(
#         description='Bulk import of emoji from a slack team'
#     )
#     parser.add_argument(
#         'directory',
#         help='Where do we store downloaded emoji?'
#     )
#     parser.add_argument(
#         '--team-name', '-t',
#         default=os.getenv('SLACK_TEAM'),
#         help='Defaults to the $SLACK_TEAM environment variable.'
#     )
#     parser.add_argument(
#         '--cookie', '-c',
#         default=os.getenv('SLACK_COOKIE'),
#         help='Defaults to the $SLACK_COOKIE environment variable.'
#     )
#     parser.add_argument(
#         '--concurrent-requests', '-r',
#         default=int(os.getenv('CONCURRENT_REQUESTS', '200')),
#         type=int,
#         help='Maximum concurrent requests. Defaults to the $CONCURRENT_REQUESTS environment variable or 200.'
#     )
#     args = parser.parse_args()
#     return args



def arg_envs(cookie:str, team_name:str, token:str)-> tuple:
    """Pull args from env if needed, and assert requirements"""
    _cookie = os.getenv("SLACK_COOKIE", cookie)
    _team_name = os.getenv("SLACK_TEAM", team_name)
    _token = os.getenv("SLACK_TOKEN", token)

    assert _cookie, "Either SLACK_COOKIE env var, or --cookie param must be set"
    assert _team_name, "Either SLACK_TEAM env var, or --team param must be set"
    assert _token, "Either SLACK_TOKEN env var, or --token param must be set"

    return (_cookie, _team_name, _token)

